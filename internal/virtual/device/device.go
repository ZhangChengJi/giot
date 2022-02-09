package device

import (
	"encoding/json"
	"fmt"
	"giot/internal/model"
	"giot/internal/virtual/mqtt"
	"time"
)

var (
	DataChan   chan *model.DeviceMsg
	AlarmChan  chan *model.DeviceMsg
	OnlineChan chan *model.DeviceMsg
)

type device struct {
	mqtt.Broker
}

type Interface interface {
	listenLoop()
	Insert(data *model.DeviceMsg)
	InsertAlarm(data *model.DeviceMsg)
	Online(data *model.DeviceMsg)
}

func Init() {
	DataChan = make(chan *model.DeviceMsg)
	AlarmChan = make(chan *model.DeviceMsg)
	OnlineChan = make(chan *model.DeviceMsg)
	d := &device{mqtt.Broker{Client: mqtt.Client}}

	for i := 0; i < 2; i++ {
		go d.listenLoop()
	}

}

func (d *device) listenLoop() {
	for {
		select {
		case data := <-DataChan:
			d.Insert(data)
			fmt.Println("我是正常数据", data.Data)
		case data := <-AlarmChan:
			d.InsertAlarm(data)
			fmt.Println("我是报警数据", data.Data)
		case data := <-OnlineChan:
			d.Online(data)
			fmt.Println("我是上下线数据", data.Type)
		case <-time.After(300 * time.Millisecond):
		}
	}
}

//发布消息
func (d *device) Insert(data *model.DeviceMsg) {
	topic := append([]byte("device/data/"), data.DeviceId...)
	payload, _ := json.Marshal(data)
	d.Publish(string(topic), payload)

}
func (d *device) InsertAlarm(data *model.DeviceMsg) {
	topic := append([]byte("device/alarm/"), data.DeviceId...)
	payload, _ := json.Marshal(data)
	d.Publish(string(topic), payload)
}
func (d *device) Online(data *model.DeviceMsg) {
	topic := append([]byte("device/online/"), data.DeviceId...)
	payload, _ := json.Marshal(data)
	d.Publish(string(topic), payload)
}
