package device

import (
	"fmt"
	"giot/internal/core/model"
	"log"
	"time"
)

var (
	DataChan  chan *model.DeviceMsg
	AlarmChan chan *model.DeviceMsg
)

type device struct {
}

type Interface interface {
	listenLoop()
	Insert(data *model.DeviceMsg)
	InsertAlarm(data *model.DeviceMsg)
	//sendNotify(data *model.DeviceMsg)
}

func Init() {
	DataChan = make(chan *model.DeviceMsg)
	AlarmChan = make(chan *model.DeviceMsg)
	d := &device{}

	for i := 0; i < 2; i++ {
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
		case <-time.After(200 * time.Millisecond):
		}
	}
}

func (d *device) Insert(data *model.DeviceMsg) {
	log.Println("我是上传善法放假啊开发开发")
	fmt.Println("我是上树")

}
func (d *device) InsertAlarm(data *model.DeviceMsg) {
	fmt.Println("我是告警")
}
