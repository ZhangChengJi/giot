package transfer

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"giot/internal/notify"
	"giot/internal/notify/sms"
	"giot/internal/scheduler/db"
	broker "giot/internal/scheduler/mqtt"
	"giot/internal/virtual/model"
	"giot/pkg/log"
	"giot/pkg/queue"
	"giot/utils"
	"giot/utils/consts"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"gorm.io/gorm"
	"runtime"
	"time"
)

var qOption *queue.Option = queue.DefaultOption().SetMaxQueueSize(10000).SetMaxBatchSize(100)

type Transfer struct {
	mqtt      *broker.Broker
	td        *sql.DB
	db        *gorm.DB
	dataChan  chan []byte
	alarmChan chan []byte
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
		alarmChan: make(chan []byte, 1024),
		notifier:  notify.New(),
		// Prepare batch queue
		q: queue.NewWithOption(qOption),
	}
	t.consume(4)
	go t.notifyLoop()
	t.ListenMqtt()
}

func (t *Transfer) ListenMqtt() {
	t.mqtt.Client.Subscribe("transfer/data/#", 0, t.dataHandler)
	t.mqtt.Client.Subscribe("transfer/alarm/#", 0, t.alarmHandler)
}
func (t *Transfer) dataHandler(client mqtt.Client, msg mqtt.Message) {
	t.q.Enqueue(msg.Payload())
	fmt.Printf("TOPIC: %s\n", msg.Topic())
	fmt.Printf("MSG: %s\n", msg.Payload())
}
func (t *Transfer) alarmHandler(client mqtt.Client, msg mqtt.Message) {
	t.alarmChan <- msg.Payload()
	t.q.Enqueue(msg.Payload())
	fmt.Printf("TOPIC: %s\n", msg.Topic())
	fmt.Printf("MSG: %s\n", msg.Payload())
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

func (t *Transfer) notifyProvider(action string, metadata *notify.Metadata, template string) {

	switch action {
	case consts.SMS:
		sms := sms.New(metadata.RegionId, metadata.AccessKeyId, metadata.AccessSecret)
		sms.AddReceivers(metadata.PhoneNumbers)
		t.notifier.UseServices(sms)
		t.notifier.Send(context.Background(), metadata.TemplateCode, template) //TODO 是否记录发送状态
	case consts.VOICE:

	}
}
func (t *Transfer) notifyLoop() {
	for {
		select {
		case alarm := <-t.alarmChan:
			var msg model.DeviceMsg
			if err := json.Unmarshal(alarm, &msg); err != nil {
				log.Errorf("json Unmarshal failed:%v", err)
				return
			}
			for _, action := range msg.Actions {
				metadata, err := t.queryNotifyData(action.NotifierId, action.TemplateId)
				if err != nil {
					log.Errorf("query notify metadata failed")
					return
				}
				template := &notify.Template{
					DeviceName: msg.Name,
					SlaveId:    msg.SlaveId,
					Value:      msg.Data,
				}
				te, _ := json.Marshal(template)

				t.notifyProvider(action.NotifyType, metadata, string(te))
			}
		case <-time.After(200 * time.Millisecond):

		}
	}
}

func (t *Transfer) queryNotifyData(cid, tid string) (*notify.Metadata, error) {
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
	err = json.Unmarshal([]byte(template.Template), &metadata)
	if err != nil {
		return nil, err
	}
	return &metadata, nil
}
