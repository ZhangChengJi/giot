package device

import (
	"context"
	"encoding/json"
	"errors"
	"giot/internal/model"
	"giot/pkg/etcd"
	"giot/pkg/log"
	"giot/pkg/modbus"
	"gorm.io/gorm"
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
	//指令下发方式 1:单条下发 2:多条下发',
	INSTRUCT_ONE  = 1
	INSTRUCT_MANY = 2
)

type Device struct {
	modbuls modbus.Client
	etcd    etcd.Interface
	db      *gorm.DB
}

func Setup(Etcd etcd.Interface, Db *gorm.DB) error {
	device := &Device{
		modbuls: modbus.NewClient(&modbus.RtuHandler{}),
		etcd:    Etcd,
		db:      Db,
	}
	return device.deviceLoad()

}

func (device *Device) deviceLoad() error {
	log.Info("Initialize device load.....")
	size := 10
	page := 1

	var devices []*model.PigDevice

	for {
		offset := size * (page - 1)
		err := device.db.Offset(offset).Limit(size).Find(&devices).Error //查询设备
		if err != nil {
			return err
		}
		//没有数据就跳出加载
		if len(devices) <= 0 {
			break
		}
		page++

		for _, d := range devices {
			var product model.PigProduct
			err := device.db.Where(&model.PigProduct{Id: d.ProductId}).First(&product).Error
			if err != nil {
				return err
			}
			var instruct bool
			ta := &model.TimerActive{Guid: d.DeviceId}
			if d.InstructFlag == INSTRUCT_ONE {
				instruct = true
			}
			slaves, err := device.getSalve(d.DeviceId, ta, instruct)
			if err != nil {
				return err
			}
			devic := &model.Device{
				GuId:        d.DeviceId,
				Name:        d.DeviceName,
				ProductType: product.ProductType,
				FCode:       ta.Ft,
				Salve:       slaves,
			}
			if !reflect.DeepEqual(devic, model.Device{}) { //etcd device
				ctx, cancel := context.WithTimeout(context.TODO(), 5*time.Second)
				data, _ := json.Marshal(devic)
				err := device.etcd.Create(ctx, "device/"+ta.Guid, string(data))
				cancel()
				if err != nil {
					return err
				}
			}
		}

	}
	return nil
}

func (device *Device) getSalve(guid string, ta *model.TimerActive, instruct bool) ([]*model.Slave, error) {

	var deviceSlaves []*model.PigDeviceSlave
	err := device.db.Where(&model.PigDeviceSlave{DeviceId: guid, SlaveStatus: DEVICE_ENABLE}).Find(&deviceSlaves).Error
	if err != nil {
		return nil, err
	}
	var slaves []*model.Slave
	var already bool
	for _, s := range deviceSlaves {
		//属性查询
		var property *model.PigProperty
		err = device.db.Where(&model.PigProperty{Id: s.PropertyId}).First(&property).Error
		if err != nil {
			log.Error("load property no found.", err)
			return nil, err
		}

		slave := &model.Slave{
			SlaveId:   byte(s.ModbusAddress),
			SlaveName: s.SlaveName,
		}
		if !reflect.DeepEqual(property, model.PigProperty{}) {

			//30s 指令
			if instruct && !already { //单
				err := device.createActive(ta, property.PropertyRegister, s.ModbusAddress, property.AddressOffset)
				if err != nil {
					return nil, err
				}
				already = true
			} else if !instruct { //多
				err := device.createActive(ta, property.PropertyRegister, s.ModbusAddress, property.AddressOffset)
				if err != nil {
					return nil, err
				}
			}

			//告警
			var deviceAlarm *model.PigPropertyAlarm
			err = device.db.Where(&model.PigPropertyAlarm{Id: property.AlarmId, AlarmStatus: DEVICE_ENABLE}).First(&deviceAlarm).Error
			if err != nil {
				log.Error("load property alarm no found.", err)
				return nil, err
			}
			//[{"level": "3", "operator": "eq", "filterValue": "30", "leftValueType": "1"}]
			//{"time": "5", "Handle": "first", "enabled": "1", "threshold": "2"}
			//告警条件
			var trigger []*model.Trigger

			err = json.Unmarshal([]byte(deviceAlarm.AlarmRule), &trigger)
			if err != nil {
				return nil, err
			}

			//防抖
			var shake *model.ShakeLimit
			err = json.Unmarshal([]byte(deviceAlarm.AlarmShake), &shake)
			if err != nil {
				return nil, err
			}
			slave.Alarm = &model.Alarm{
				Triggers:   trigger,
				ShakeLimit: shake,
			}

		}
		slaves = append(slaves, slave)
	}
	return slaves, nil
}

// createActive  创建动态指令
func (device *Device) createActive(ta *model.TimerActive, code, salveId int, propertyRegister int) error {
	tc := &tCode{}
	err := tc.crateTimerCode(TIMER30_SECOND, ta).functionCode(code, device.modbuls, byte(salveId), uint16(propertyRegister))
	if err != nil {
		return err
	}
	return err

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
	if ta.Ft == nil {
		ta.Ft = &model.Ft{Tm: duration}
	}

	if ta.Ft.Tm == duration {
		tc.Ft = ta.Ft
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
