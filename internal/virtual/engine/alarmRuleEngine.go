package engine

import (
	model2 "giot/internal/model"
	"giot/internal/virtual/device"
	"giot/utils/consts"
	"time"
)

type Interface interface {
	AlarmRule(unit uint16, slave *model2.Slave)
	Trigger(data uint16, slave *model2.Slave)
	Action(guid, name, productId, alarmId string, alarmLevel int, data uint16, slaveId byte, actions []*model2.Action)
}

type AlarmRuleEngine struct {
	Alarms []*model2.Alarm
}

func NewAlarmRule(alarms []*model2.Alarm) *AlarmRuleEngine {
	return &AlarmRuleEngine{
		Alarms: alarms,
	}
}

func (engine *AlarmRuleEngine) AlarmRule(data uint16, slave *model2.Slave) {
	engine.Trigger(data, slave)
}

func (engine *AlarmRuleEngine) Trigger(data uint16, slave *model2.Slave) {
	device.DataChan <- &model2.DeviceMsg{
		Ts:        time.Now(),
		Type:      consts.DATA,
		Status:    true,
		DeviceId:  slave.DeviceId,
		SlaveId:   int(slave.SlaveId),
		Name:      slave.SlaveName,
		ProductId: slave.ProductId,
		Data:      data,
	}
	for _, alarm := range engine.Alarms { //循环告警规则
	loop:
		for _, trigger := range alarm.Triggers { //循环告警触发条件
			if trigger.Type == "properties" { //判断是否是属性
				if trigger.ModelId == slave.AttributeId { //判断是否当前属性ID
					switch trigger.Operator { //判断比对条件(任意) 触发条件满足条件中任意一个即可触发
					case consts.EQ: //=
						if data == trigger.Val {
							engine.Action(slave.DeviceId, alarm.DeviceName, alarm.ProductId, alarm.AlarmId, alarm.AlarmLevel, data, slave.SlaveId, alarm.Actions)
							break loop
						}
					case consts.NOT:
						if data != trigger.Val {
							engine.Action(slave.DeviceId, alarm.DeviceName, alarm.ProductId, alarm.AlarmId, alarm.AlarmLevel, data, slave.SlaveId, alarm.Actions)
							break loop
						}
					case consts.GT:
						if data > trigger.Val {
							engine.Action(slave.DeviceId, alarm.DeviceName, alarm.ProductId, alarm.AlarmId, alarm.AlarmLevel, data, slave.SlaveId, alarm.Actions)
							break loop
						}
					case consts.LT:
						if data < trigger.Val {
							engine.Action(slave.DeviceId, alarm.DeviceName, alarm.ProductId, alarm.AlarmId, alarm.AlarmLevel, data, slave.SlaveId, alarm.Actions)
							break loop
						}
					case consts.GTE:
						if data >= trigger.Val {
							engine.Action(slave.DeviceId, alarm.DeviceName, alarm.ProductId, alarm.AlarmId, alarm.AlarmLevel, data, slave.SlaveId, alarm.Actions)
							break loop
						}
					case consts.LTE:
						if data <= trigger.Val {
							engine.Action(slave.DeviceId, alarm.DeviceName, alarm.ProductId, alarm.AlarmId, alarm.AlarmLevel, data, slave.SlaveId, alarm.Actions)
							break loop
						}
					}
				}
			}
		}
	}
}
func (engine *AlarmRuleEngine) Action(guid, name, productId, alarmId string, alarmLevel int, data uint16, slaveId byte, actions []*model2.Action) {
	device.AlarmChan <- &model2.DeviceMsg{
		Ts:         time.Now(),
		Type:       consts.ALARM,
		Status:     false,
		DeviceId:   guid,
		SlaveId:    int(slaveId),
		Name:       name,
		ProductId:  productId,
		AlarmId:    alarmId,
		AlarmLevel: alarmLevel,
		Data:       data,
		Actions:    actions,
	}
}
