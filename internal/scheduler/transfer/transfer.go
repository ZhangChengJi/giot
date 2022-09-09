package transfer

import (
	"database/sql"
	"giot/internal/scheduler/line"
	"giot/internal/virtual/device"
	"giot/pkg/log"
	"giot/pkg/queue"
	"giot/utils"
	"giot/utils/consts"
	"giot/utils/json"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/go-redis/redis"
	"golang.org/x/text/message"
	"gorm.io/gorm"
	"runtime"
	"strconv"
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
	li        line.LineCache
}

func Setup(mqtt mqtt.Client, tdengine *sql.DB, mysql *gorm.DB, redis *redis.Client) {
	var t = &Transfer{
		mq:        mqtt,
		td:        tdengine,
		db:        mysql,
		li:        &line.Line{Re: redis},
		alarmChan: make(chan *device.DeviceMsg, 1024),
		queue:     queue.NewWithOption(queue.DefaultOption()),
	}
	t.consume(t.queue, 4, t.td)
	t.listenMqtt()
}

func (t *Transfer) listenMqtt() {
	t.mq.Subscribe("device/data/#", 0, t.dataHandler)
	t.mq.Subscribe("device/alarm/#", 0, t.alarmHandler)
	t.mq.Subscribe("device/line/#", 0, t.lineHandler)
}
func (t *Transfer) dataHandler(client mqtt.Client, msg mqtt.Message) {
	var device device.DeviceMsg
	if err := FromMqttBytes(msg.Payload(), &device); err != nil {
		log.Sugar.Errorf("topic data:%v Err: %v\n", msg.Topic(), message.Key, err)
		return
	}
	if device.GroupId > 0 { //判断是否分组id未空，没有分配分组就不进行数据存储
		t.queue.Enqueue(device)
	}

	//fmt.Printf("TOPIC: %s\n", msg.Topic())
	//fmt.Printf("MSG: %s\n", msg.Payload())
}
func (t *Transfer) alarmHandler(client mqtt.Client, msg mqtt.Message) {
	var device device.DeviceMsg
	if err := FromMqttBytes(msg.Payload(), &device); err != nil {
		log.Sugar.Errorf("topic data:%v Err: %v\n", msg.Topic(), message.Key, err)
		return
	}
	if device.GroupId > 0 { //判断是否分组id未空，没有分配分组就不进行数据存储
		t.alarmChan <- &device
		t.queue.Enqueue(device)
	}
	//fmt.Printf("TOPIC: %s\n", msg.Topic())
	//fmt.Printf("MSG: %s\n", msg.Payload())
}

func (t *Transfer) lineHandler(client mqtt.Client, msg mqtt.Message) {
	var device device.DeviceMsg

	if err := FromMqttBytes(msg.Payload(), &device); err != nil {
		log.Sugar.Errorf("topic data:%v Err: %v\n", msg.Topic(), message.Key, err)
		return
	}
	switch device.Status {
	case consts.ONLINE:
		t.online(device.DeviceId, device.SlaveId)
		log.Sugar.Warnf("guid:%s online", device.DeviceId)
	case consts.OFFLINE:
		t.offline(device.DeviceId, device.SlaveId)
		log.Sugar.Warnf("guid:%s offline", device.DeviceId)
	}

}

//设备↕上下线
func (t *Transfer) online(deviceId string, slaveId int) {
	if slaveId == 0 {
		t.li.SetDeviceOnline(deviceId)
	} else {
		t.li.SetSlaveOnline(deviceId, strconv.Itoa(slaveId))
	}

}
func (t *Transfer) offline(deviceId string, slaveId int) {
	if slaveId == 0 {
		t.li.SetDeviceOffline(deviceId)
		t.li.BatchSlaveOffline(deviceId)
	} else {
		t.li.SetSlaveOffline(deviceId, strconv.Itoa(slaveId))
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
					log.Sugar.Errorf("cannot build sql with records: %v", err)
					continue
				}
				runtime.Gosched()
				//fmt.Println("sql:======" + sqlStr)
				_, err = taos.Exec(sqlStr)
				if err != nil {
					log.Sugar.Errorf("exec query error: %v, the sql command is:\n%s\n", err, sqlStr)
				}

			}
		}(q)
	}
}
