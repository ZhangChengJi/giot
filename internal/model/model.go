package model

import (
	"fmt"
	"giot/internal/virtual/device"
	"giot/utils/consts"
	"giot/utils/encoding"
	"github.com/panjf2000/gnet/v2"
	"github.com/shopspring/decimal"
	"strconv"
	"strings"
	"time"
)

type RemoteData struct {
	Frame      []byte
	RemoteAddr string
	Conn       gnet.Conn
}
type RegisterData struct {
	Conn gnet.Conn
	D    string
}
type PigProduct struct {
	Id             int    `gorm:"id" json:"id"`
	ProductType    int    `gorm:"product_type" json:"product_type"`       // 产品类型 1:工业 2:家报
	ProductName    string `gorm:"product_name" json:"product_name"`       // 产品名称
	ProductFactory string `gorm:"product_factory" json:"product_factory"` // 设备厂家
	ProductModel   string `gorm:"product_model" json:"product_model"`     // 设备型号
	ProductImg     string `gorm:"product_img" json:"product_img"`         // 产品图片
	ProductDesc    string `gorm:"product_desc" json:"product_desc"`       // 设备描述

}

func (PigProduct) TableName() string {

	return "pig_product"
}

type PigDevice struct {
	Id               string `gorm:"id" json:"id"`                               // 设备id
	ProductId        int    `gorm:"product_id" json:"product_id"`               // 产品id
	DeviceName       string `gorm:"device_name" json:"device_name"`             // 设备名称
	NetworkFlag      int    `gorm:"network_flag" json:"network_flag"`           // 联网方式 1:4G-DTU 1:NB-IOT
	InstructFlag     int    `gorm:"instruct_flag" json:"instruct_flag"`         // 指令下发方式 1:单条下发 2:多条下发
	SimCode          string `gorm:"sim_code" json:"sim_code"`                   // SIM卡ICCID
	BindStatus       int    `gorm:"bind_status" json:"bind_status"`             // 绑定状态 0:未绑定 1:已绑定
	LineStatus       int    `gorm:"line_status" json:"line_status"`             // 设备状态: 0:离线 1:在线
	DeviceAddress    string `gorm:"device_address" json:"device_address"`       // 设备地址
	DeviceCoordinate string `gorm:"device_coordinate" json:"device_coordinate"` // 设备坐标信息
	GroupId          int    `gorm:"group_id" json:"group_id"`                   //分组id

}

func (PigDevice) TableName() string {
	return "pig_device"
}

type PigDeviceSlave struct {
	DeviceId      string `gorm:"device_id" json:"device_id"`           // 设备ID
	SlaveAlias    string `gorm:"slave_alias" json:"slave_alias"`       // 从机设备别名
	SlaveName     string `gorm:"slave_name" json:"slave_name"`         // 从机设备名称
	ModbusAddress int    `gorm:"modbus_address" json:"modbus_address"` // modbus从站地址
	PropertyId    int    `gorm:"property_id" json:"property_id"`       // 关联设备属性
	SlaveDesc     string `gorm:"slave_desc" json:"slave_desc"`         // 从机设备描述
	SlaveStatus   int    `gorm:"slave_status" json:"slave_status"`     // 从机设备开关 0:关闭 1:开启
	LineStatus    int    `gorm:"line_status" json:"line_status"`       // 从机设备状态 0:离线 1:在线
	CreateTime    string `gorm:"create_time" json:"create_time"`       // 创建时间
	UpdateTime    string `gorm:"update_time" json:"update_time"`       // 修改时间
	CreateBy      string `gorm:"create_by" json:"create_by"`           // 创建者
	UpdateBy      string `gorm:"update_by" json:"update_by"`           // 更新人
}

func (PigDeviceSlave) TableName() string {
	return "pig_device_slave"
}

type PigProperty struct {
	Id                     int    `gorm:"id" json:"id"`
	GroupId                int    `gorm:"group_id" json:"group_id"`                               // 分类id
	GroupName              string `gorm:"group_name" json:"group_name"`                           // 分类名称
	AlarmId                int    `gorm:"alarm_id" json:"alarm_id"`                               // 告警id
	PropertyName           string `gorm:"property_name" json:"property_name"`                     // 属性名称
	PropertyIdentification string `gorm:"property_identification" json:"property_identification"` // 属性标识
	PropertyDataType       string `gorm:"property_data_type" json:"property_data_type"`           // 数据类型
	PropertyPrecision      int    `gorm:"property_precision" json:"property_precision"`           //浮点型精度
	PropertyUnit           string `gorm:"property_unit" json:"property_unit"`                     // 单位
	PropertyRegister       int    `gorm:"property_register" json:"property_register"`             // 寄存器
	AddressOffset          int    `gorm:"address_offset" json:"address_offset"`                   // 地址偏移
	PropertyImg            string `gorm:"property_img" json:"property_img"`                       // 属性图标
	PropertyDesc           string `gorm:"property_desc" json:"property_desc"`                     // 属性描述
}

func (PigProperty) TableName() string {
	return "pig_property"
}

type PigPropertyAlarm struct {
	Id          int    `gorm:"id" json:"id"`
	AlarmStatus int    `gorm:"alarm_status" json:"alarm_status"` // 报警状态 0:关闭 1:开启
	AlarmLevel  string `gorm:"alarm_level" json:"alarm_level"`   // 告警等级:1:低级2:中级3:高级
	AlarmRule   string `gorm:"alarm_rule" json:"alarm_rule"`     // 告警条件
	AlarmShake  string `gorm:"alarm_shake" json:"alarm_shake"`   // 告警防抖
	AlarmNotice string `gorm:"alarm_notice" json:"alarm_notice"` // 告警通知

}

func (PigPropertyAlarm) TableName() string {
	return "pig_property_alarm"
}

type Trigger struct {
	FilterValue   uint16 `json:"filterValue"`   //过滤值
	LeftValueType string `json:"leftValueType"` //左侧值
	Level         int    `json:"level"`         //报警等级
	Operator      string `json:"operator"`      //比对条件
}

type ShakeLimit struct { //防抖
	Enabled   int    `json:"enabled"`   //是否开启防抖
	Time      int    `json:"time"`      //时间间隔(秒)
	Threshold int    `json:"threshold"` //触发阈值(次)
	Handle    string `json:"Handle"`    //是否第一次满足条件就触发

}

// Device 结构体
type Device struct {
	GuId         string   `json:"guid"`
	Name         string   `json:"name"`
	ProductType  int      `json:"productType"`
	ProductModel string   `json:"productModel"`
	BindStatus   int      `json:"bindStatus"`
	Instruct     int      `json:"instruct"`
	LineStatus   int      `json:"lineStatus"`
	GroupId      int      `json:"groupId"`
	FCode        *Ft      `json:"fCode"`
	Salve        []*Slave `json:"salve"`
}

func (device *Device) IsType() bool {

	if device.ProductType == 1 { //工业
		return true
	} else {
		return false
	}
}
func (device *Device) IsInstruct() bool {

	if device.Instruct == 1 { //是否是单指令
		return true
	} else {
		return false
	}
}

type ListenMsg struct {
	ListenType int
	RemoteAddr string
	Guid       string
	Command    Comm
}
type TimerActive struct {
	Guid string `json:"guId"`
	Ft   *Ft    `json:"ft"`
}
type Ft struct {
	Tm    time.Duration `json:"tm"`
	FCode [][]byte      `json:"fCode"`
}

type Slave struct {
	SlaveId    byte      `json:"slaveId"`
	SlaveName  string    `json:"slaveName"`
	Alarm      *Alarm    `json:"alarm"`
	DataTime   time.Time `json:"dataTime"`
	LineStatus string    `json:"lineStatus"`
	Precision  int       `json:"precision"`
}

type Comm int8

// '告警级别（1--低，2--中，3--高）'
func Level(a int) string {
	switch a {
	case 1:
		return "低级"
	case 2:
		return "中级"
	case 3:
		return "高级"
	}
	return ""
}

type DeviceChange struct {
	ChangeType string `json:"changeType"`
	DeviceId   string `json:"deviceId"`
}

type Interface interface {
	AlarmRule(slaveId byte, point int, data []byte, fcode uint8, info *Device)
	Trigger(slaveId byte, point int, data []byte, info *Device)
	Action(guid, status string, level int, data float64, slaveId byte, groupId int)
}

type Alarm struct {
	Triggers   []*Trigger  `json:"trigger"`    //条件
	ShakeLimit *ShakeLimit `json:"shakeLimit"` //防抖动配置
}

func (engine *Alarm) AlarmRule(slaveId byte, point int, data []byte, fcode uint8, info *Device) {
	value := encoding.BytesToUint16(encoding.BIG_ENDIAN, data)
	if info.IsType() {
		switch value { //是否工业产品
		// 10000     探测器内部错误
		case consts.InternalError:
			engine.Action(info.GuId, consts.ALARM, consts.Internal, 0, slaveId, info.GroupId)
			break

		// 20000     通讯错误
		case consts.CommunicationError:
			engine.Action(info.GuId, consts.ALARM, consts.Communication, 0, slaveId, info.GroupId)
			break

		// 30000     主机未连接探测器、主机屏蔽探测器
		case consts.ShieldError:
			engine.Action(info.GuId, consts.ALARM, consts.Shield, 0, slaveId, info.GroupId)
			break

			// 65535     探头故障
		case consts.SlaveHitchError:
			engine.Action(info.GuId, consts.ALARM, consts.SlaveHitch, 0, slaveId, info.GroupId)
			break
		}
		engine.Trigger(slaveId, point, data, info)
	} else {
		switch fcode {
		case consts.ReadCode:
			engine.Trigger(slaveId, point, data, info)
			break
		case consts.HomeHitchError:
			//故障（若为0，是传感器低故障报警，若为1，是传感器高故障报警，若为2，是传感器寿命报警）
			if value == 0 { //若为0，是传感器低故障报警
				engine.Action(info.GuId, consts.ALARM, consts.LowHitch, 0, slaveId, info.GroupId)
			} else if value == 1 { //若为1，是传感器高故障报警
				engine.Action(info.GuId, consts.ALARM, consts.HighHitch, 0, slaveId, info.GroupId)
			} else if value == 2 { //若为2，是传感器寿命报警
				engine.Action(info.GuId, consts.ALARM, consts.Life, 0, slaveId, info.GroupId)
			}
			break
		case consts.HomeHighError:
			engine.Action(info.GuId, consts.ALARM, consts.High, 0, slaveId, info.GroupId)
			break
		case consts.HomeLowError:
			engine.Action(info.GuId, consts.ALARM, consts.Low, 0, slaveId, info.GroupId)
			break
		}

	}
}
func floating(point int, data []byte) (val string) {
	str := strconv.Itoa(int(encoding.BytesToUint16(encoding.BIG_ENDIAN, data)))
	spot := len(str) - point
	if len(str) > point {
		for i, v := range str {
			if i < spot {
				val = strings.Join([]string{val, string(v)}, "")
			} else if i == spot {
				val = strings.Join([]string{val, string(v)}, ".")
			} else {
				val = strings.Join([]string{val, string(v)}, "")
			}
		}
	}
	return
}
func (engine *Alarm) Trigger(slaveId byte, point int, data []byte, info *Device) {

	if point > 0 {
		str := floating(point, data)
		fmt.Printf("没有转---->string:%v", data)
		dec, err := decimal.NewFromString(str)
		fmt.Printf("转成---->decimal:%v", data)
		value, _ := dec.Float64()
		fmt.Printf("转成---->Float64:%v", value)
		if err != nil {
			return
		}
		for _, trigger := range engine.Triggers { //循环告警触发条件
			filter := decimal.NewFromInt32(int32(trigger.FilterValue))
			switch trigger.Operator { //TODO 判断比对条件(任意) 触发条件满足条件中任意一个即可触发  高报优先
			case consts.EQ: //=

				if dec.Equal(filter) {
					engine.Action(info.GuId, consts.ALARM, trigger.Level, value, slaveId, info.GroupId)
					return
				}
			case consts.NOT:
				if !dec.Equal(filter) {
					engine.Action(info.GuId, consts.ALARM, trigger.Level, value, slaveId, info.GroupId)
					return
				}
			case consts.GT:
				if dec.GreaterThan(filter) {
					engine.Action(info.GuId, consts.ALARM, trigger.Level, value, slaveId, info.GroupId)
					return
				}
			case consts.LT:
				if dec.LessThan(filter) {
					engine.Action(info.GuId, consts.ALARM, trigger.Level, value, slaveId, info.GroupId)
					return
				}
			case consts.GTE:
				if dec.GreaterThanOrEqual(filter) {
					engine.Action(info.GuId, consts.ALARM, trigger.Level, value, slaveId, info.GroupId)
					return
				}
			case consts.LTE:
				if dec.LessThanOrEqual(filter) {
					engine.Action(info.GuId, consts.ALARM, trigger.Level, value, slaveId, info.GroupId)
					return
				}
			}
		}

		engine.Action(info.GuId, consts.DATA, consts.Normal, value, slaveId, info.GroupId)

	} else {
		data := encoding.BytesToUint16(encoding.BIG_ENDIAN, data)
		in := decimal.NewFromInt32(int32(data))
		value, _ := in.Float64()
		for _, trigger := range engine.Triggers { //循环告警触发条件
			switch trigger.Operator { //TODO 判断比对条件(任意) 触发条件满足条件中任意一个即可触发  高报优先
			case consts.EQ: //=
				if data == trigger.FilterValue {
					engine.Action(info.GuId, consts.ALARM, trigger.Level, value, slaveId, info.GroupId)
					return
				}
			case consts.NOT:
				if data != trigger.FilterValue {
					engine.Action(info.GuId, consts.ALARM, trigger.Level, value, slaveId, info.GroupId)
					return
				}
			case consts.GT:
				if data > trigger.FilterValue {
					engine.Action(info.GuId, consts.ALARM, trigger.Level, value, slaveId, info.GroupId)
					return
				}
			case consts.LT:
				if data < trigger.FilterValue {
					engine.Action(info.GuId, consts.ALARM, trigger.Level, value, slaveId, info.GroupId)
					return
				}
			case consts.GTE:
				if data >= trigger.FilterValue {
					engine.Action(info.GuId, consts.ALARM, trigger.Level, value, slaveId, info.GroupId)
					return
				}
			case consts.LTE:
				if data <= trigger.FilterValue {
					engine.Action(info.GuId, consts.ALARM, trigger.Level, value, slaveId, info.GroupId)
					return
				}
			}
		}

		engine.Action(info.GuId, consts.DATA, consts.Normal, value, slaveId, info.GroupId)
	}

}

func (engine *Alarm) Action(guid, status string, level int, data float64, slaveId byte, groupId int) {

	device.DataChan <- &device.DeviceMsg{
		Ts:       time.Now(),
		DataType: consts.DATA,
		Level:    level,
		DeviceId: guid,
		SlaveId:  int(slaveId),
		Data:     data,
		GroupId:  groupId,
	}
	if status == consts.ALARM {
		device.AlarmChan <- &device.DeviceMsg{
			Ts:       time.Now(),
			DataType: consts.ALARM,
			Level:    level,
			DeviceId: guid,
			SlaveId:  int(slaveId),
			Data:     data,
			GroupId:  groupId,
		}
	}

}
