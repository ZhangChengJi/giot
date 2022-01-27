package engine

import (
	"giot/internal/virtual/device"
	"giot/internal/virtual/model"
	"giot/utils/consts"
)

var EnChan = make(chan model.RemoteData, 1024)

//func loop() {
//
//	data := <-EnChan
//}

type Interface interface {
	AlarmRule(unit float64, slave *model.Slave)
	Trigger(data float64, slave *model.Slave)
	Action(guid, name, productId string, data float64, actions []*model.Action)
}

type AlarmRuleEngine struct {
	Alarms []*model.Alarm
}

func NewAlarmRule(alarms []*model.Alarm) *AlarmRuleEngine {
	return &AlarmRuleEngine{
		Alarms: alarms,
	}
}

func (engine *AlarmRuleEngine) AlarmRule(data float64, slave *model.Slave) {
	engine.Trigger(data, slave)
}

func (engine *AlarmRuleEngine) Trigger(data float64, slave *model.Slave) {
	var b bool
	for _, alarm := range engine.Alarms { //循环告警规则
		for _, trigger := range alarm.Triggers { //循环告警触发条件
			if trigger.Type == "properties" { //判断是否是属性
				if trigger.ModelId == slave.AttributeId { //判断是否当前属性ID
					switch trigger.Operator { //判断比对条件(任意)
					case consts.EQ, consts.NOT, consts.GT, consts.LT, consts.GTE, consts.LTE: //=
						if data == trigger.Val {
							engine.Action(slave.DeviceId, alarm.Name, alarm.ProductId, data, alarm.Actions)
							b = true
							break
						}

					}

				}
			}
		}
		if !b {
			device.DataChan <- &model.DeviceMsg{
				Type:      consts.DATA,
				Status:    true,
				DeviceId:  slave.DeviceId,
				Name:      slave.SlaveName,
				ProductId: slave.ProductId,
				Data:      data,
			}
		}
	}
}
func (engine *AlarmRuleEngine) Action(guid, name, productId string, data float64, actions []*model.Action) {
	device.DataChan <- &model.DeviceMsg{
		Type:      consts.DATA,
		Status:    false,
		DeviceId:  guid,
		Name:      name,
		ProductId: productId,
		Data:      data,
	}
	device.AlarmChan <- &model.DeviceMsg{
		Type:      consts.ALARM,
		Status:    false,
		DeviceId:  guid,
		Name:      name,
		ProductId: productId,
		Data:      data,
		Actions:   actions,
	}
}
