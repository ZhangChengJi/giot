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
	"giot/utils/consts"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"gorm.io/gorm"
)

type Transfer struct {
	mqtt      *broker.Broker
	td        *sql.DB
	db        *gorm.DB
	dataChan  chan []byte
	alarmChan chan []byte
	notifier  *notify.Notify
}

func SetupTransfer() {
	t := &Transfer{
		mqtt: &broker.Broker{
			Client: broker.Client,
		},
		td:        db.Td,
		db:        db.DB,
		dataChan:  make(chan []byte, 1024),
		alarmChan: make(chan []byte, 1024),
		notifier:  notify.New(),
	}
	go t.loop()
	t.ListenMqtt()
}

func (t *Transfer) ListenMqtt() {
	t.mqtt.Client.Subscribe("transfer/data/#", 0, func(client mqtt.Client, msg mqtt.Message) {
		t.dataChan <- msg.Payload()
		fmt.Printf("TOPIC: %s\n", msg.Topic())
		fmt.Printf("MSG: %s\n", msg.Payload())
	})
	t.mqtt.Client.Subscribe("transfer/alarm/#", 0, func(client mqtt.Client, msg mqtt.Message) {
		t.alarmChan <- msg.Payload()
		//t.insert()
		notify.Provider("sms")
		fmt.Printf("TOPIC: %s\n", msg.Topic())
		fmt.Printf("MSG: %s\n", msg.Payload())
	})
}

func (t *Transfer) insert(msg *model.DeviceMsg) {

	//t.db.Exec()
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
func (t *Transfer) loop() {
	for {
		select {
		case data := <-t.dataChan:
			var msg model.DeviceMsg
			if err := json.Unmarshal(data, &msg); err != nil {
				log.Errorf("json Unmarshal failed:%v", err)
				return
			}
			t.insert(&msg)
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
					SlaveName:  msg.SlaveName,
					Value:      msg.Data,
				}
				te, _ := json.Marshal(template)

				t.notifyProvider(action.NotifyType, metadata, string(te))
			}

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
