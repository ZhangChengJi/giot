package device

import (
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
	//sendNotify(data *model.DeviceMsg)
}

func Init() {
	DataChan = make(chan *model.DeviceMsg)
	AlarmChan = make(chan *model.DeviceMsg)
	d := &device{mqtt.Broker{Client: mqtt.Client}}

	for i := 0; i < 5; i++ {
		go d.listenLoop()
	}

}

func (d *device) listenLoop() {
	for {
		select {
		case data := <-DataChan:
			d.Insert(data)
		case data := <-AlarmChan:
			d.InsertAlarm(data)
		case data := <-OnlineChan:
			d.Online(data)
		case <-time.After(200 * time.Millisecond):
		}
	}
}

//发布消息
func (d *device) Insert(data *model.DeviceMsg) {
	topic := append([]byte("transfer/data/"), data.DeviceId...)
	topic = append(topic, ""...)
	d.Publish(string(topic), data)

}
func (d *device) InsertAlarm(data *model.DeviceMsg) {
	topic := append([]byte("transfer/alarm/"), data.DeviceId...)
	d.Publish(string(topic), data)
}
func (d *device) Online(data *model.DeviceMsg) {
	topic := append([]byte("transfer/online/"), data.DeviceId...)
	d.Publish(string(topic), data)
}
