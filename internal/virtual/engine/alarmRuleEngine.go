package engine

import (
	"giot/internal/virtual/device"
	"giot/internal/virtual/model"
	"giot/pkg/modbus"
	"giot/utils"
	"giot/utils/consts"
	"strconv"
)

var EnChan = make(chan model.RemoteData, 1024)

//func loop() {
//
//	data := <-EnChan
//}

type Interface interface {
	AlarmRule(unit *modbus.ProtocolDataUnit, slave *model.Slave)
	Trigger(data []byte, slave *model.Slave)
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

func (engine *AlarmRuleEngine) AlarmRule(unit *modbus.ProtocolDataUnit, slave *model.Slave) {
	engine.Trigger(unit.Data, slave)
}

func (engine *AlarmRuleEngine) Trigger(data []byte, slave *model.Slave) {
	for _, alarm := range engine.Alarms { //循环告警规则
		for _, trigger := range alarm.Triggers { //循环告警触发条件
			if trigger.Type == "properties" { //判断是否是属性
				if trigger.ModelId == slave.AttributeId { //判断是否当前属性ID
					switch trigger.Operator { //判断比对条件
					case "eq": //=

						as, _ := strconv.ParseUint(string(data), 16, 32)
						as1, _ := strconv.ParseUint(string(trigger.Val), 16, 32)
						if as == as1 {
							engine.Action(slave.DeviceId, alarm.Name, alarm.ProductId, utils.ByteToFloat64(data), alarm.Actions)
							break
						}
					case "not": //<>
						if utils.ByteToFloat64(data) != utils.ByteToFloat64(trigger.Val) {
							engine.Action(slave.DeviceId, alarm.Name, alarm.ProductId, utils.ByteToFloat64(data), alarm.Actions)
							break
						}
					case "gt": //>
						if utils.ByteToFloat64(data) > utils.ByteToFloat64(trigger.Val) {
							engine.Action(slave.DeviceId, alarm.Name, alarm.ProductId, utils.ByteToFloat64(data), alarm.Actions)
							break
						}
					case "lt": //<
						if utils.ByteToFloat64(data) < utils.ByteToFloat64(trigger.Val) {
							engine.Action(slave.DeviceId, alarm.Name, alarm.ProductId, utils.ByteToFloat64(data), alarm.Actions)
							break
						}
					case "gte": //>=
						if utils.ByteToFloat64(data) >= utils.ByteToFloat64(trigger.Val) {
							engine.Action(slave.DeviceId, alarm.Name, alarm.ProductId, utils.ByteToFloat64(data), alarm.Actions)
							break
						}
					case "lte": //<=
						if utils.ByteToFloat64(data) <= utils.ByteToFloat64(trigger.Val) {
							engine.Action(slave.DeviceId, alarm.Name, alarm.ProductId, utils.ByteToFloat64(data), alarm.Actions)
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
		Type:      consts.ALARM,
		DeviceId:  guid,
		Name:      name,
		ProductId: productId,
		Data:      data,
		Actions:   actions,
	}
}
