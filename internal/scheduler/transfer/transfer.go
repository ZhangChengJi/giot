package transfer

import (
	"context"
	"database/sql"
	"fmt"
	"giot/internal/virtual/device"
	"giot/utils"
	"giot/utils/json"
	"golang.org/x/text/message"
	"runtime"
	"strconv"

	"giot/internal/model"
	"giot/internal/notify"
	"giot/internal/notify/sms"
	"giot/pkg/log"
	"giot/pkg/queue"
	"giot/utils/consts"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"gorm.io/gorm"
	"time"
)

var qOption *queue.Option = queue.DefaultOption().SetMaxQueueSize(10000).SetMaxBatchSize(100)

type Transfer struct {
	mq        mqtt.Client
	td        *sql.DB
	db        *gorm.DB
	dataChan  chan []byte
	alarmChan chan *device.DeviceMsg
	notifier  *notify.Notify
	queue     *queue.Queue
}

func Setup(mqtt mqtt.Client, tdengine *sql.DB, mysql *gorm.DB) {
	var t = &Transfer{
		mq:        mqtt,
		td:        tdengine,
		db:        mysql,
		alarmChan: make(chan *device.DeviceMsg, 1024),
		notifier:  notify.New(),
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

func (t *Transfer) notifyProvider(action string, metadata *notify.Metadata) {

	switch action {
	case consts.SMS:
		sms := sms.New(metadata.AccessKeyId, metadata.Secret)
		sms.AddReceivers(metadata.Sms.PhoneNumber)
		t.notifier.UseServices(sms)
		t.notifier.Send(context.Background(), metadata.Sms.SignName, metadata.Sms.Code, metadata.Sms.Param) //TODO 是否记录发送状态
	case consts.VOICE:

	}
}

//func (t *Transfer) notifyLoop() {
//	for {
//		select {
//		case alarm := <-t.alarmChan:
//			for _, action := range alarm.Actions {
//				metadata, err := t.queryNotifyMetadata(action.NotifierId, action.TemplateId, action.NotifyType, alarm.Name, alarm.SlaveId, alarm.AlarmLevel)
//				if err != nil {
//					log.Errorf("query notify metadata failed")
//					return
//				}
//				t.notifyProvider(action.NotifyType, metadata)
//			}
//		case <-time.After(300 * time.Millisecond):
//
//		}
//	}
//}

func (t *Transfer) queryNotifyMetadata(cid, tid, notifyType, name string, slaveId int, level int) (*notify.Metadata, error) {
	var config model.NotifyConfig
	err := t.db.First(&config, cid).Error
	if err != nil {
		return nil, err
	}
	var metadata notify.Metadata
	err = json.Unmarshal([]byte(config.Configuration), &metadata)
	if err != nil {
		return nil, err
	}
	var template model.NotifyTemplate
	err = t.db.First(&template, tid).Error
	if err != nil {
		return nil, err
	}
	switch notifyType {
	case consts.SMS:
		pa := &sms.Template{Devname: name, Devid: strconv.Itoa(slaveId), Alarmtype: model.Level(level)}
		param, _ := json.Marshal(pa)
		err = json.Unmarshal([]byte(template.Template), &metadata.Sms)
		if err != nil {
			return nil, err
		}
		metadata.Sms.Param = string(param)
		break
	case consts.VOICE:
		err = json.Unmarshal([]byte(template.Template), &metadata.Voice)
		if err != nil {
			return nil, err
		}
		break
	}

	return &metadata, nil
}
