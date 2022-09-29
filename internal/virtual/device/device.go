package device

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"giot/internal/virtual/mqtt"
	"giot/pkg/log"
	"giot/pkg/modbus"
	mqtts "github.com/eclipse/paho.mqtt.golang"
	"strings"
	"time"
)

var (
	DataChan      chan *DeviceMsg
	LastChan      chan *DeviceMsg
	AlarmChan     chan *DeviceMsg
	AlarmSaveChan chan *DeviceMsg
	NotifyChan    chan *DeviceMsg
	OnlineChan    chan *DeviceMsg
	DebugChan     chan *Debug
)

type Debug struct {
	Guid  string `json:"guid"`
	FCode []byte `json:"fCode"`
}
type device struct {
	mqtt   mqtt.Broker
	modbus modbus.Client
}

type Interface interface {
	listenLoop()
	Insert(data *DeviceMsg)
	InsertAlarm(data *DeviceMsg)
	Line(data *DeviceMsg)
}

func Init() {
	DataChan = make(chan *DeviceMsg, 1024)
	AlarmChan = make(chan *DeviceMsg, 1024)
	AlarmSaveChan = make(chan *DeviceMsg, 1024)
	NotifyChan = make(chan *DeviceMsg, 1024)
	OnlineChan = make(chan *DeviceMsg, 1024)
	LastChan = make(chan *DeviceMsg, 10240)
	DebugChan = make(chan *Debug)
	d := &device{mqtt: mqtt.Broker{Client: mqtt.Client}, modbus: modbus.NewClient(&modbus.RtuHandler{})}

	for i := 0; i < 5; i++ {
		go d.listenLoop()
	}
	go d.Subscribe()

}

func (d *device) listenLoop() {
	for {
		select {
		case data := <-LastChan: //实时最新数据
			d.InsertLast(data)
		case data := <-DataChan: //存储数据
			d.Insert(data)
		case data := <-AlarmChan: //告警实时数据
			d.InsertAlarm(data)
		case data := <-AlarmSaveChan: //告警实时数据
			d.InsertAlarmSave(data)
		case data := <-OnlineChan: //上下线数据
			d.Line(data)
		case data := <-NotifyChan: //通知数据
			d.InsertNotify(data)

		case <-time.After(200 * time.Millisecond):
		}
	}
}

//发布消息
func (d *device) Insert(data *DeviceMsg) {
	var buf bytes.Buffer
	buf.WriteString("device/data/")
	buf.WriteString(data.DeviceId)
	payload, _ := json.Marshal(data)
	d.mqtt.Publish(buf.String(), payload)

}
func (d *device) InsertLast(data *DeviceMsg) {
	var buf bytes.Buffer
	buf.WriteString("device/last/")
	buf.WriteString(data.DeviceId)
	payload, _ := json.Marshal(data)
	d.mqtt.Publish(buf.String(), payload)

}
func (d *device) InsertNotify(data *DeviceMsg) {
	if data.GroupId > 0 {
		topic := append([]byte("device/notify/"), data.DeviceId...)
		payload, _ := json.Marshal(data)
		d.mqtt.Publish(string(topic), payload)
	}
}
func (d *device) InsertAlarm(data *DeviceMsg) {
	topic := append([]byte("device/alarm/"), data.DeviceId...)
	payload, _ := json.Marshal(data)
	d.mqtt.Publish(string(topic), payload)
}
func (d *device) InsertAlarmSave(data *DeviceMsg) {
	topic := append([]byte("device/alarmSave/"), data.DeviceId...)
	payload, _ := json.Marshal(data)
	d.mqtt.Publish(string(topic), payload)
}
func (d *device) Line(data *DeviceMsg) {
	topic := append([]byte("device/line/"), data.DeviceId...)
	payload, _ := json.Marshal(data)
	d.mqtt.Publish(string(topic), payload)
}

func (d *device) Subscribe() {
	d.mqtt.Client.Subscribe("device/debug/#", 0, d.Debug)
}
func (d *device) Debug(client mqtts.Client, msg mqtts.Message) {
	if len(msg.Topic()) > 5 && len(msg.Topic()) < 40 && len(msg.Payload()) < 20 && len(msg.Payload()) > 0 {
		if len(msg.Payload()) == 16 {
			//string数据      string转hex
			data := strings.Split(msg.Topic(), "/")
			if len(msg.Payload()) > 0 {
				dst := strings.Trim(string(msg.Payload()), " ")
				if len(dst) > 0 {
					pl, err := hex.DecodeString(dst)
					if err == nil {
						if err := d.modbus.CheckCrc(pl); err == nil {
							log.Sugar.Infof("debug 调试 topic:%v", msg.Topic())
							if len(data) == 3 {
								DebugChan <- &Debug{
									Guid:  data[2],
									FCode: pl,
								}
							}
						} else {
							log.Sugar.Errorf("crc解析错误，下发数据是否是hex格式？")
						}

					} else {
						log.Sugar.Errorf("转移hex错误：%v", err)
					}
				}
			}
		} else if len(msg.Payload()) == 8 {
			//hex数据

			data := strings.Split(msg.Topic(), "/")
			if err := d.modbus.CheckCrc(msg.Payload()); err == nil {
				log.Sugar.Infof("debug 调试 topic:%v", msg.Topic())
				if len(data) == 3 {
					DebugChan <- &Debug{
						Guid:  data[2],
						FCode: msg.Payload(),
					}
				}
			} else {
				log.Sugar.Errorf("crc解析错误，下发数据是否是hex格式？")
			}
		}
	}
}
