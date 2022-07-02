package device

import (
	"bytes"
	"encoding/json"
	"fmt"
	"giot/internal/virtual/mqtt"
	"time"
)

var (
	DataChan   chan *DeviceMsg
	AlarmChan  chan *DeviceMsg
	OnlineChan chan *DeviceMsg
)

type device struct {
	mqtt.Broker
}

type Interface interface {
	listenLoop()
	Insert(data *DeviceMsg)
	InsertAlarm(data *DeviceMsg)
	Online(data *DeviceMsg)
}

func Init() {
	DataChan = make(chan *DeviceMsg)
	AlarmChan = make(chan *DeviceMsg)
	OnlineChan = make(chan *DeviceMsg)
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
			fmt.Println("我是上下线数据", data.Status)
		case <-time.After(200 * time.Millisecond):
		}
	}
}

//发布消息
func (d *device) Insert(data *DeviceMsg) {
	var buf bytes.Buffer
	buf.WriteString("device/data/")
	buf.WriteString(data.DeviceId)
	fmt.Println(buf.String())
	payload, _ := json.Marshal(data)
	d.Publish(buf.String(), payload)

}
func (d *device) InsertAlarm(data *DeviceMsg) {
	topic := append([]byte("device/alarm/"), data.DeviceId...)
	payload, _ := json.Marshal(data)
	d.Publish(string(topic), payload)
}
func (d *device) Online(data *DeviceMsg) {
	topic := append([]byte("device/online/"), data.DeviceId...)
	payload, _ := json.Marshal(data)
	d.Publish(string(topic), payload)
}
