package device

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"giot/internal/model"
	"giot/internal/scheduler/line"
	"giot/pkg/etcd"
	"giot/pkg/log"
	"giot/pkg/modbus"
	"giot/utils/consts"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/go-redis/redis"
	"github.com/schollz/progressbar/v3"
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
	TIMER1_SECOND  = 1 * time.Second
	TIMER2_SECOND  = 2 * time.Second
	TIMER3_SECOND  = 3 * time.Second
	TIMER4_SECOND  = 4 * time.Second
	TIMER5_SECOND  = 5 * time.Second
	TIMER6_SECOND  = 6 * time.Second
	TIMER7_SECOND  = 7 * time.Second
	TIMER8_SECOND  = 8 * time.Second
	TIMER9_SECOND  = 9 * time.Second
	TIMER10_SECOND = 10 * time.Second
	TIMER20_SECOND = 20 * time.Second
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
	mqtt    mqtt.Client
	li      line.LineCache
}

func Setup(Etcd etcd.Interface, Db *gorm.DB, mqtt mqtt.Client, redis *redis.Client) error {
	device := &Device{
		modbuls: modbus.NewClient(&modbus.RtuHandler{}),
		etcd:    Etcd,
		db:      Db,
		mqtt:    mqtt,
		li:      &line.Line{Re: redis},
	}

	err := device.deviceLoad()
	if err != nil {
		return err
	}
	go device.DeviceLister()
	return err

}
func (device *Device) clearAllOnline() {
	device.li.ClearAll()
}

//监听设备发生配置变化
func (device *Device) DeviceLister() error {
	if token := device.mqtt.Subscribe("device/change", 0, func(client mqtt.Client, message mqtt.Message) {
		var changeData model.DeviceChange
		err := json.Unmarshal(message.Payload(), &changeData)
		if err != nil {
			log.Sugar.Errorf("topic：'device/change' deviceData failed:%s", err)
			return
		}
		if changeData.ChangeType == consts.Update || changeData.ChangeType == consts.Add {
			var de model.PigDevice
			err := device.db.Where(&model.PigDevice{Id: changeData.DeviceId}).First(&de).Error
			if err != nil {
				return
			}
			var product model.PigProduct
			err = device.db.Where(&model.PigProduct{Id: de.ProductId}).First(&product).Error
			if err != nil {
				return
			}
			var instruct bool
			ta := &model.TimerActive{Guid: de.Id}
			if de.InstructFlag == INSTRUCT_ONE {
				instruct = true
			}
			slaves, err := device.getSalve(de.Id, ta, instruct)
			if err != nil {
				return
			}
			devic := &model.Device{
				GuId:        de.Id,
				Name:        de.DeviceName,
				ProductType: product.ProductType,
				FCode:       ta.Ft,
				Salve:       slaves,
				BindStatus:  de.BindStatus,
			}
			if !reflect.DeepEqual(devic, model.Device{}) { //etcd device
				ctx, cancel := context.WithTimeout(context.TODO(), 5*time.Second)
				data, _ := json.Marshal(devic)
				err := device.etcd.Create(ctx, "device/"+ta.Guid, string(data))
				cancel()
				if err != nil {
					return
				}
			}

		}
		if changeData.ChangeType == consts.Delete {
			if err = device.etcd.Delete(context.TODO(), "device/"+changeData.DeviceId); err != nil {
				log.Sugar.Errorf("delete etcd deivceId :%v", changeData.DeviceId)
			}

		}
	}); token.Wait() && token.Error() != nil {
		fmt.Println(token.Error())
		log.Sugar.Errorf("subscribe topic:%s failed", "device/change")
		return token.Error()

	}
	return nil
}

func (device *Device) deviceLoad() error {
	log.Sugar.Info("Initialize device load.....")
	log.Sugar.Info("清理之前的在线状态")
	device.clearAllOnline()
	size := 100
	page := 1
	log.Sugar.Info("删除之前的元数据")
	device.etcd.DeleteWithPrefix(context.TODO(), "device/")
	var devices []*model.PigDevice
	var cou int64
	device.db.Model(devices).Count(&cou)
	bar := progressbar.Default(cou)

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
			bar.Clear()
			bar.Add(1)
			var product model.PigProduct
			err := device.db.Where(&model.PigProduct{Id: d.ProductId}).First(&product).Error
			if err != nil {
				return err
			}
			var instruct bool
			ta := &model.TimerActive{Guid: d.Id}
			if d.InstructFlag == INSTRUCT_ONE {
				instruct = true
			}
			slaves, err := device.getSalve(d.Id, ta, instruct)
			if err != nil {
				log.Sugar.Errorf("load db slave info error：%s", err)
				continue
			}
			devic := &model.Device{
				GuId:        d.Id,
				Name:        d.DeviceName,
				ProductType: product.ProductType,
				Instruct:    d.InstructFlag,
				FCode:       ta.Ft,
				Salve:       slaves,
				BindStatus:  d.BindStatus,
				LineStatus:  d.LineStatus,
				GroupId:     d.GroupId,
				Address:     d.DeviceAddress,
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
	//30s 指令
	if instruct && !already { //单
		err := device.createF1Active(ta, 3, len(deviceSlaves))
		if err != nil {
			return nil, err
		}
		already = true
	}
	for _, s := range deviceSlaves {
		//属性查询
		var property *model.PigProperty
		err = device.db.Where(&model.PigProperty{Id: s.PropertyId}).First(&property).Error
		if err != nil {
			log.Sugar.Errorf("load property no found.", err)
			return nil, err
		}

		slave := &model.Slave{
			SlaveId:      byte(s.ModbusAddress),
			SlaveName:    s.SlaveName,
			Precision:    property.PropertyPrecision,
			PropertyUnit: property.PropertyUnit,
			PropertyName: property.PropertyName,
		}
		if !reflect.DeepEqual(property, model.PigProperty{}) {

			if !instruct { //多
				err := device.createActive(ta, property.PropertyRegister, s.ModbusAddress, property.AddressOffset)
				if err != nil {
					return nil, err
				}
			}

			//告警
			var deviceAlarm *model.PigPropertyAlarm
			err = device.db.Where(&model.PigPropertyAlarm{Id: property.AlarmId, AlarmStatus: DEVICE_ENABLE}).First(&deviceAlarm).Error
			if err != nil {
				log.Sugar.Error("load property alarm no found.", err)
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

			slaves = append(slaves, slave)
		}
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
func (device *Device) createF1Active(ta *model.TimerActive, code, salveSize int) error {
	tc := &tCode{}
	err := tc.crateTimerCode(TIMER30_SECOND, ta).functionCode(code, device.modbuls, byte(241), uint16(salveSize))
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
