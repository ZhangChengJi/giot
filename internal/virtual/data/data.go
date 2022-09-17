package data

import (
	"context"
	"errors"
	"giot/conf"
	"giot/internal/model"
	"giot/pkg/log"
	"giot/pkg/modbus"
	"giot/utils/json"
	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"
	"time"

	gorm1 "giot/pkg/gorm"
	redis1 "giot/pkg/redis"
)

const (
	DEVICE_INFO    = "device_info"
	SLAVE_INFO     = "slave_info"
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

type Data struct {
	modbuls modbus.Client
	db      *gorm.DB
	re      *redis.Client
}

func New() *Data {
	mysql, err := gorm1.New(conf.MysqlConfig)
	if err != nil {
		log.Sugar.Errorf("mysql 加载配置错误%s", err)
		return nil
	}
	redi, err := redis1.New(conf.RedisConfig)
	if err != nil {
		log.Sugar.Errorf("redis 加载配置错误%s", err)
		return nil
	}
	log.Sugar.Info("redis mysql 加载完成...")
	return &Data{
		modbuls: modbus.NewClient(&modbus.RtuHandler{}),
		db:      mysql,
		re:      redi,
	}
}
func (d *Data) GetTimerData(deviceId string) (active *model.TimerActive, err error) {
	device, err := d.GetDevice(deviceId)
	if err != nil {
		return
	}
	salve, err := d.getSalve(deviceId)
	if err != nil {
		return
	}
	var instruct bool
	if device.Instruct == INSTRUCT_ONE {
		instruct = true
	}
	active, err = d.getActive(deviceId, instruct, salve)
	if err != nil {
		return
	}
	return active, nil

}
func (d *Data) GetSlaveData(deviceId string) (slaves []*model.Slave, err error) {
	slaves, err = d.getSalve(deviceId)
	if err != nil {
		return nil, err
	}
	return slaves, err
}

func (d *Data) GetData(deviceId string) (devic *model.Device, err error) {

	device, err := d.GetDevice(deviceId)
	slaves, err := d.getSalve(deviceId)
	var instruct bool
	if device.Instruct == INSTRUCT_ONE {
		instruct = true
	}
	active, err := d.getActive(deviceId, instruct, slaves)
	devic = &model.Device{
		GuId:        device.GuId,
		Name:        device.Name,
		ProductType: 1,
		Instruct:    device.Instruct,
		FCode:       active.Ft,
		Salve:       slaves,
		BindStatus:  device.BindStatus,
		LineStatus:  device.LineStatus,
		GroupId:     device.GroupId,
		Address:     device.Address,
	}
	return devic, err

}
func (d *Data) GetDevice(deviceId string) (devic *model.Device, err error) {
	k, err := d.re.Get(context.TODO(), DEVICE_INFO+":"+deviceId).Bytes()

	if len(k) > 0 {
		err = json.Unmarshal(k, &devic)
		if err != nil {
			log.Sugar.Errorf("转换错误%v错误", deviceId)
			return nil, err
		}
		return devic, err
	} else {
		var device *model.PigDevice
		err := d.db.Where(&model.PigDevice{Id: deviceId}).First(&device).Error
		if err != nil {
			log.Sugar.Errorf("缓存失败")
			return nil, err
		}
		devic = &model.Device{
			GuId:        device.Id,
			Name:        device.DeviceName,
			ProductType: 1,
			Instruct:    device.InstructFlag,
			BindStatus:  device.BindStatus,
			LineStatus:  device.LineStatus,
			GroupId:     device.GroupId,
			Address:     device.DeviceAddress,
		}
		dejson, err := json.Marshal(devic)
		if err != nil {
			return nil, err
		}
		err = d.re.Set(context.TODO(), DEVICE_INFO+":"+device.Id, dejson, 0).Err()
		return devic, err
	}
}
func (d *Data) getActive(deviceId string, instruct bool, slaves []*model.Slave) (active *model.TimerActive, err error) {
	active = &model.TimerActive{Guid: deviceId}
	if instruct {
		err := d.createF1Active(active, 3, len(slaves))
		if err != nil {
			return nil, err
		}
	} else {
		for _, slave := range slaves {
			err := d.createActive(active, slave.PropertyRegister, int(slave.SlaveId), slave.AddressOffset)
			if err != nil {
				return nil, err
			}
		}
	}
	return active, err
}
func (d *Data) getSalve(guid string) (slaves []*model.Slave, err error) {
	result, err := d.re.Get(context.TODO(), SLAVE_INFO+":"+guid).Bytes()

	if len(result) > 0 {
		err = json.Unmarshal(result, &slaves)
		if err != nil {
			return nil, err
		}
		return slaves, nil
	}
	var deviceSlaves []*model.PigDeviceSlave
	err = d.db.Model(&model.PigDeviceSlave{}).Select(
		"pig_device_slave.slave_name,"+
			"pig_device_slave.modbus_address,"+
			"p.property_identification,"+
			"p.property_unit,"+
			"p.property_name,"+
			"p.property_register,"+
			"p.address_offset,"+
			"r.alarm_rule ").Joins("LEFT JOIN pig_property p ON p.id=pig_device_slave.property_id "+
		"LEFT JOIN pig_device_rule r ON r.property_id =p.id").Where(
		"pig_device_slave.device_id=? AND r.device_id =? ", guid, guid).Scan(&deviceSlaves).Error
	if err != nil {
		log.Sugar.Errorf("load slave property rule no found.", err)
		return nil, err
	}

	for _, s := range deviceSlaves {

		slave := &model.Slave{
			SlaveId:          byte(s.ModbusAddress),
			SlaveName:        s.SlaveName,
			Precision:        s.PropertyPrecision,
			PropertyUnit:     s.PropertyUnit,
			PropertyName:     s.PropertyName,
			PropertyRegister: s.PropertyRegister,
			AddressOffset:    s.AddressOffset,
		}
		//告警条件
		var trigger []*model.Trigger

		err = json.Unmarshal([]byte(s.AlarmRule), &trigger)
		if err != nil {
			return nil, err
		}

		slave.Rule = &model.Rule{
			Triggers: trigger,
		}

		slaves = append(slaves, slave)
	}
	sljson, err := json.Marshal(slaves)
	if err != nil {
		return nil, err
	}
	err = d.re.Set(context.TODO(), SLAVE_INFO+":"+guid, sljson, 0).Err()
	if err != nil {
		log.Sugar.Errorf("缓存失败")
		return nil, err
	}
	return slaves, nil
}

// createActive  创建动态指令
func (d *Data) createActive(ta *model.TimerActive, code, salveId int, propertyRegister int) error {
	tc := &tCode{}
	err := tc.crateTimerCode(TIMER30_SECOND, ta).functionCode(code, d.modbuls, byte(salveId), uint16(propertyRegister))
	if err != nil {
		return err
	}
	return err

}
func (d *Data) createF1Active(ta *model.TimerActive, code, salveSize int) error {
	tc := &tCode{}
	err := tc.crateTimerCode(TIMER30_SECOND, ta).functionCode(code, d.modbuls, byte(241), uint16(salveSize))
	if err != nil {
		return err
	}
	return err

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
