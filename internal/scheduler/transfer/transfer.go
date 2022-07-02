package transfer

import (
	"database/sql"
	"fmt"
	"giot/internal/model"
	"giot/internal/virtual/device"
	"giot/pkg/log"
	"giot/pkg/queue"
	"giot/utils"
	"giot/utils/consts"
	"giot/utils/json"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"golang.org/x/text/message"
	"gorm.io/gorm"
	"runtime"
	"time"
)

var qOption *queue.Option = queue.DefaultOption().SetMaxQueueSize(10000).SetMaxBatchSize(100)

type Transfer struct {
	mq        mqtt.Client
	td        *sql.DB
	db        *gorm.DB
	dataChan  chan []byte
	alarmChan chan *device.DeviceMsg
	queue     *queue.Queue
}

func Setup(mqtt mqtt.Client, tdengine *sql.DB, mysql *gorm.DB) {
	var t = &Transfer{
		mq:        mqtt,
		td:        tdengine,
		db:        mysql,
		alarmChan: make(chan *device.DeviceMsg, 1024),
		queue:     queue.NewWithOption(queue.DefaultOption()),
	}
	t.consume(t.queue, 4, t.td)
	//go t.notifyLoop()
	t.listenMqtt()
}

func (t *Transfer) listenMqtt() {
	t.mq.Subscribe("device/data/#", 0, t.dataHandler)
	t.mq.Subscribe("device/alarm/#", 0, t.alarmHandler)
	t.mq.Subscribe("device/online/#", 0, t.onlineHandler)
}
func (t *Transfer) dataHandler(client mqtt.Client, msg mqtt.Message) {
	var device device.DeviceMsg
	if err := FromMqttBytes(msg.Payload(), &device); err != nil {
		log.Errorf("topic data:%v Err: %v\n", msg.Topic(), message.Key, err)
		return
	}
	t.queue.Enqueue(device)
	fmt.Printf("TOPIC: %s\n", msg.Topic())
	fmt.Printf("MSG: %s\n", msg.Payload())
}
func (t *Transfer) alarmHandler(client mqtt.Client, msg mqtt.Message) {
	var device device.DeviceMsg
	if err := FromMqttBytes(msg.Payload(), &device); err != nil {
		log.Errorf("topic data:%v Err: %v\n", msg.Topic(), message.Key, err)
		return
	}
	t.alarmChan <- &device
	t.queue.Enqueue(device)
	fmt.Printf("TOPIC: %s\n", msg.Topic())
	fmt.Printf("MSG: %s\n", msg.Payload())
}

func (t *Transfer) onlineHandler(client mqtt.Client, msg mqtt.Message) {
	var device device.DeviceMsg
	if err := FromMqttBytes(msg.Payload(), &device); err != nil {
		log.Errorf("topic data:%v Err: %v\n", msg.Topic(), message.Key, err)
		return
	}
	switch device.Status {
	case consts.ONLINE:
		t.online(device.DeviceId, device.SlaveId)
		log.Warnf("guid:%s online", device.DeviceId)
	case consts.OFFLINE:
		t.offline(device.DeviceId, device.SlaveId)
		log.Warnf("guid:%s offline", device.DeviceId)
	}

}

//设备↕上下线
func (t *Transfer) online(guid string, slaveId int) {
	if slaveId == 0 {
		var pigDevice model.PigDevice
		err := t.db.Debug().Model(&pigDevice).Where("line_status=? and id=?", 0, guid).Update("line_status", 1).Error
		if err != nil {
			log.Errorf("online guid:%s update failed", guid)
			return
		}

	} else {
		var slave model.PigDeviceSlave
		fmt.Printf("slave:%v online", slaveId)
		err := t.db.Debug().Model(&slave).Where("device_id=? and modbus_address=? and line_status=? ", guid, slaveId, 0).Update("line_status", 1).Error
		if err != nil {
			log.Errorf("online guid:%s update failed", guid)
			return
		}
	}

}
func (t *Transfer) offline(guid string, slaveId int) {
	var slave model.PigDeviceSlave
	if slaveId == 0 {
		var pigDevice model.PigDevice
		err := t.db.Debug().Model(&pigDevice).Where("line_status=? and id=?", 1, guid).Update("line_status", 0).Error
		if err != nil {
			log.Errorf("offline guid:%s update failed", guid)
			return
		}
		//fmt.Printf("slave:%v offline", slaveId)
		//err = t.db.Debug().Model(&slave).Where("id=? and line_status=? ", guid, 1).Update("line_status", 0).Error
		//if err != nil {
		//	log.Errorf("online guid:%s update failed", guid)
		//	return
		//}
	} else {
		fmt.Printf("slave:%v offline", slaveId)
		err := t.db.Debug().Model(&slave).Where("device_id=? and modbus_address=? and line_status=? ", guid, slaveId, 1).Update("line_status", 0).Error
		if err != nil {
			log.Errorf("online guid:%s update failed", guid)
			return
		}
	}
}

func FromMqttBytes(bytes []byte, device *device.DeviceMsg) error {
	return json.Unmarshal(bytes, &device)
}

func (t *Transfer) consume(q *queue.Queue, workers int, taos *sql.DB) {
	for i := 0; i < workers; i++ {
		go func(q *queue.Queue) {
			for {
				runtime.Gosched()
				msg, err := q.Dequeue()
				if err != nil {
					time.Sleep(200 * time.Millisecond)
					continue
				}
				sqlStr, err := utils.ToTaosBatchInsertSql(msg)
				if err != nil {
					log.Errorf("cannot build sql with records: %v", err)
					continue
				}
				runtime.Gosched()
				fmt.Println("sql:======" + sqlStr)
				_, err = taos.Exec(sqlStr)
				if err != nil {
					log.Errorf("exec query error: %v, the sql command is:\n%s\n", err, sqlStr)
				}

			}
		}(q)
	}
}
