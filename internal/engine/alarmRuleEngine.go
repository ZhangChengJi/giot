package engine

import (
	"giot/internal/core/model"
	"giot/internal/device"
	"giot/internal/manager/modbus"
	"giot/utils"
	"strconv"
)

var EnChan = make(chan model.RemoteData, 1024)

//func loop() {
//
//	data := <-EnChan
//}

type Interface interface {
	AlarmRule(unit *modbus.ProtocolDataUnit, attributeId string, guid string)
	Trigger(attributeId string, data []byte, guid string)
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

func (engine *AlarmRuleEngine) AlarmRule(unit *modbus.ProtocolDataUnit, attributeId string, guid string) {
	engine.Trigger(attributeId, unit.Data, guid)
}

func (engine *AlarmRuleEngine) Trigger(attributeId string, data []byte, guid string) {
	for _, alarm := range engine.Alarms { //循环告警规则
		for _, trigger := range alarm.Triggers { //循环告警触发条件
			if trigger.Type == "properties" { //判断是否是属性
				if trigger.ModelId == attributeId { //判断是否当前属性ID
					switch trigger.Operator { //判断比对条件
					case "eq": //=

						as, _ := strconv.ParseUint(string(data), 16, 32)
						as1, _ := strconv.ParseUint(string(trigger.Val), 16, 32)
						if as == as1 {
							engine.Action(guid, alarm.Name, alarm.ProductId, utils.ByteToFloat64(data), alarm.Actions)
							break
						}
					case "not": //<>
						if utils.ByteToFloat64(data) != utils.ByteToFloat64(trigger.Val) {
							engine.Action(guid, alarm.Name, alarm.ProductId, utils.ByteToFloat64(data), alarm.Actions)
							break
						}
					case "gt": //>
						if utils.ByteToFloat64(data) > utils.ByteToFloat64(trigger.Val) {
							engine.Action(guid, alarm.Name, alarm.ProductId, utils.ByteToFloat64(data), alarm.Actions)
							break
						}
					case "lt": //<
						if utils.ByteToFloat64(data) < utils.ByteToFloat64(trigger.Val) {
							engine.Action(guid, alarm.Name, alarm.ProductId, utils.ByteToFloat64(data), alarm.Actions)
							break
						}
					case "gte": //>=
						if utils.ByteToFloat64(data) >= utils.ByteToFloat64(trigger.Val) {
							engine.Action(guid, alarm.Name, alarm.ProductId, utils.ByteToFloat64(data), alarm.Actions)
							break
						}
					case "lte": //<=
						if utils.ByteToFloat64(data) <= utils.ByteToFloat64(trigger.Val) {
							engine.Action(guid, alarm.Name, alarm.ProductId, utils.ByteToFloat64(data), alarm.Actions)
							break
						}

					}

				}
			}
		}
	}
}
func (engine *AlarmRuleEngine) Action(guid, name, productId string, data float64, actions []*model.Action) {
	device.AlarmChan <- &model.DeviceMsg{
		DeviceId:  guid,
		Name:      name,
		ProductId: productId,
		Data:      data,
		Actions:   actions,
	}
}
