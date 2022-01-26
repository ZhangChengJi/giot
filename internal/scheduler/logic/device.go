package logic

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"giot/internal/scheduler/db"
	"giot/internal/scheduler/model"
	"giot/pkg/etcd"
	"giot/pkg/modbus"
	"reflect"
	"time"
)

var (
	DEVICE_DISABLE = 0 //关闭
	DEVICE_ENABLE  = 1 //启用
	DEVICE_ACTIVE  = 1 //主动上报
	DEVICE_PASSIVE = 2 //服务器采集
	F03H           = 3 //03功能
	F04H           = 4 //04功能

	TIMER30_SECOND = 30 * time.Second
	TIMER60_SECOND = 60 * time.Second
)

type DeviceSvc struct {
	Modbus modbus.Client
	Etcd   etcd.Interface
}

func (device *DeviceSvc) InitEtcdDataLoad() error {
	size := 2
	page := 1
	var devices []model.Device
	for {
		offset := size * (page - 1)
		err := db.DB.Offset(offset).Limit(size).Where(&model.Device{EnableStatus: DEVICE_ENABLE}).Find(&devices).Error //查询设备
		if err != nil {
			return err
		}
		//没有数据就跳出加载
		if len(devices) <= 0 {
			break
		}
		page++

		for _, d := range devices {

			var detectors []model.Detector
			err = db.DB.Where(&model.Detector{DeviceId: d.Id, SlaveDeviceSwitch: DEVICE_ENABLE}).Find(&detectors).Error //查询从机
			if err != nil {
				return err
			}
			var ta = &model.TimerActive{Guid: d.Id} //定时指令struct
			var savle []*model.Slave

			for _, de := range detectors {
				var attribute model.Attribute
				db.DB.Where(&model.Attribute{Id: de.AttributeId}).First(&attribute) //查询属性
				if !reflect.DeepEqual(attribute, model.Attribute{}) {               //属性不为空
					fmt.Println(attribute)
					if attribute.SourceType == DEVICE_PASSIVE { //是否为服务采集
						tc := &tCode{}
						switch attribute.Frequency { //判断采集频率
						case 1: //30秒
							tc.crateTimerCode(TIMER30_SECOND, ta).functionCode(attribute.FunctionCode, device.Modbus, byte(de.SlaveAddress), uint16(attribute.RegAddress))
						case 2: //60秒
							tc.crateTimerCode(TIMER60_SECOND, ta).functionCode(attribute.FunctionCode, device.Modbus, byte(de.SlaveAddress), uint16(attribute.RegAddress))

						}
					}
					sa := &model.Slave{
						ProductId:   d.ProductId,
						DeviceName:  d.Name,
						DeviceId:    d.Id,
						SlaveId:     byte(de.SlaveAddress),
						SlaveName:   de.Name,
						AttributeId: attribute.Id,
					}
					savle = append(savle, sa)
				}
			}
			//****************告警规则********************
			var alarms []*model.AlarmRule
			db.DB.Where(&model.AlarmRule{ProductId: d.ProductId, EnableStatus: DEVICE_ENABLE}).Find(&alarms)
			var alarmList []*model.Alarm
			for _, alarm := range alarms {
				isShake := false
				isFirst := false
				if alarm.IsShake == 1 {
					isShake = true
				}
				if alarm.IsFirst == 1 {
					isFirst = false
				}

				//***************触发条件****************
				var condition []*model.Condition
				db.DB.Where(&model.Condition{AlarmRuleId: alarm.Id}).Find(&condition)
				var triggers []*model.Trigger
				for _, c := range condition {
					trigger := &model.Trigger{
						Type:     c.ConditionType,
						ModelId:  c.AttributeId,
						Operator: c.SymbolType,
						Val:      []byte(c.Values),
					}
					triggers = append(triggers, trigger)
				}
				//*******************END****************

				//**************执行动作*****************
				var executeActions []*model.ExecuteAction
				db.DB.Where(&model.ExecuteAction{AlarmRuleId: alarm.Id}).Find(&executeActions)
				var actions []*model.Action
				for _, e := range executeActions {
					actions = append(actions, &model.Action{Type: e.Type, NotifyType: e.NotifyType, NotifierId: e.NotifyConfigId, TemplateId: e.NotifyTemplateId})
				}
				//*****************END******************

				var t = &model.Alarm{
					ProductId:  d.ProductId,
					DeviceId:   d.Id,
					DeviceName: d.Name,
					ShakeLimit: &model.ShakeLimit{Enabled: isShake, Time: alarm.Within, Threshold: alarm.Num, AlarmFirst: isFirst},
					Triggers:   triggers,
					Actions:    actions,
				}
				alarmList = append(alarmList, t)
			}
			//****************END**********************

			if len(ta.Ft) > 0 { //etcd 指令存储
				ctx, cancel := context.WithTimeout(context.TODO(), 5*time.Second)
				data, _ := json.Marshal(ta)
				err := device.Etcd.Create(ctx, "transfer/"+ta.Guid+"/code", string(data))
				cancel()
				if err != nil {
					return err
				}

			}
			if len(savle) > 0 { //etcd 从机+属性
				ctx, cancel := context.WithTimeout(context.TODO(), 5*time.Second)
				data, _ := json.Marshal(savle)
				err := device.Etcd.Create(ctx, "transfer/"+ta.Guid+"/salve", string(data))
				cancel()
				if err != nil {
					return err
				}
			}
			if len(alarms) > 0 { //etcd 告警规则
				ctx, cancel := context.WithTimeout(context.TODO(), 5*time.Second)
				data, _ := json.Marshal(alarmList)
				err := device.Etcd.Create(ctx, "transfer/"+ta.Guid+"/alarm", string(data))
				cancel()
				if err != nil {
					return err
				}
			}

		}

	}

	return nil
}
func Operator(vale string) {
	const (
		eq   = "="
		not  = "<>"
		gt   = ">"
		lt   = "<"
		gte  = ">="
		lte  = "<="
		like = "like"
	)

}

type tCode struct {
	*model.Ft
}

//判断是否有同类定时数据
func (tc *tCode) crateTimerCode(duration time.Duration, ta *model.TimerActive) *tCode {
	if len(ta.Ft) <= 0 {
		ta.Ft = append(ta.Ft, &model.Ft{Tm: duration})
	}
	for i, v := range ta.Ft {
		if v.Tm == duration {
			tc.Ft = v
			break
		}
		if i == len(ta.Ft) {
			ta.Ft = append(ta.Ft, &model.Ft{Tm: duration})
		}
	}

	return tc

}

func (tc *tCode) functionCode(t int, m modbus.Client, salveId byte, regAddr uint16) error {
	switch t { //判断功能码
	case F03H: //03功能码
		code, _ := m.ReadHoldingRegisters(salveId, 0, regAddr)
		tc.FCode = append(tc.FCode, code)
		break
	case F04H:
		code, _ := m.ReadInputRegisters(salveId, 0, regAddr)
		tc.FCode = append(tc.FCode, code)
		break
	default:
		//没有对应的功能码不添加
		return errors.New("not function code")
	}
	return nil
}
