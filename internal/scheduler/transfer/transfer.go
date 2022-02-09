package transfer

import (
	"context"
	"database/sql"
	"fmt"
	"giot/internal/scheduler/logic"
	"giot/utils"
	"giot/utils/json"
	"golang.org/x/text/message"
	"runtime"
	"strconv"

	"giot/internal/model"
	"giot/internal/notify"
	"giot/internal/notify/sms"
	"giot/internal/scheduler/db"
	broker "giot/internal/scheduler/mqtt"
	"giot/pkg/log"
	"giot/pkg/queue"
	"giot/utils/consts"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"gorm.io/gorm"
	"time"
)

var qOption *queue.Option = queue.DefaultOption().SetMaxQueueSize(10000).SetMaxBatchSize(100)

type Transfer struct {
	mqtt      *broker.Broker
	td        *sql.DB
	db        *gorm.DB
	dataChan  chan []byte
	alarmChan chan *model.DeviceMsg
	notifier  *notify.Notify
	q         *queue.Queue
}

func SetupTransfer() {
	t := &Transfer{
		mqtt: &broker.Broker{
			Client: broker.Client,
		},
		td:        db.Td,
		db:        db.DB,
		alarmChan: make(chan *model.DeviceMsg, 1024),
		notifier:  notify.New(),
		// Prepare batch queue
		q: queue.NewWithOption(qOption),
	}
	t.consume(4)
	go t.notifyLoop()
	t.ListenMqtt()
}

func (t *Transfer) ListenMqtt() {
	t.mqtt.Client.Subscribe("device/data/#", 0, t.dataHandler)
	t.mqtt.Client.Subscribe("device/alarm/#", 0, t.alarmHandler)
	t.mqtt.Client.Subscribe("device/online/#", 0, t.onlineHandler)
}
func (t *Transfer) dataHandler(client mqtt.Client, msg mqtt.Message) {
	var device model.DeviceMsg
	if err := FromMqttBytes(msg.Payload(), &device); err != nil {
		log.Errorf("topic data:%v Err: %v\n", msg.Topic(), message.Key, err)
		return
	}
	t.q.Enqueue(device)
	fmt.Printf("TOPIC: %s\n", msg.Topic())
	fmt.Printf("MSG: %s\n", msg.Payload())
}
func (t *Transfer) alarmHandler(client mqtt.Client, msg mqtt.Message) {
	var device model.DeviceMsg
	if err := FromMqttBytes(msg.Payload(), &device); err != nil {
		log.Errorf("topic data:%v Err: %v\n", msg.Topic(), message.Key, err)
		return
	}

	t.alarmChan <- &device
	t.q.Enqueue(device)
	fmt.Printf("TOPIC: %s\n", msg.Topic())
	fmt.Printf("MSG: %s\n", msg.Payload())
}

func (t *Transfer) onlineHandler(client mqtt.Client, msg mqtt.Message) {
	var device model.DeviceMsg
	if err := FromMqttBytes(msg.Payload(), &device); err != nil {
		log.Errorf("topic data:%v Err: %v\n", msg.Topic(), message.Key, err)
		return
	}
	switch device.Type {
	case consts.ONLINE:
		logic.Online(device.DeviceId)
		log.Warnf("guid:%s online", device.DeviceId)
	case consts.OFFLINE:
		logic.Offline(device.DeviceId)
		log.Warnf("guid:%s offline", device.DeviceId)
	}

}
func FromMqttBytes(bytes []byte, device *model.DeviceMsg) error {
	return json.Unmarshal(bytes, &device)
}

func (t *Transfer) consume(workers int) {
	for i := 0; i < workers; i++ {
		go func(q *queue.Queue, taos *sql.DB) {
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
				_, err = taos.Exec(sqlStr)
				if err != nil {
					log.Errorf("exec query error: %v, the sql command is:\n%s\n", err, sqlStr)
				}

			}
		}(t.q, t.td)
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
func (t *Transfer) notifyLoop() {
	for {
		select {
		case alarm := <-t.alarmChan:
			for _, action := range alarm.Actions {
				metadata, err := t.queryNotifyMetadata(action.NotifierId, action.TemplateId, action.NotifyType, alarm.Name, alarm.SlaveId, alarm.AlarmLevel)
				if err != nil {
					log.Errorf("query notify metadata failed")
					return
				}
				t.notifyProvider(action.NotifyType, metadata)
			}
		case <-time.After(300 * time.Millisecond):

		}
	}
}

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
