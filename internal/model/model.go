package model

import (
	"giot/internal/virtual/device"
	"giot/utils/consts"
	"github.com/panjf2000/gnet/v2"
	"time"
)

type RemoteData struct {
	Frame      []byte
	RemoteAddr string
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
	DeviceId         string `gorm:"device_id" json:"device_id"`                 // 设备id
	ProductId        int    `gorm:"product_id" json:"product_id"`               // 产品id
	DeviceName       string `gorm:"device_name" json:"device_name"`             // 设备名称
	NetworkFlag      int    `gorm:"network_flag" json:"network_flag"`           // 联网方式 1:4G-DTU 1:NB-IOT
	InstructFlag     int    `gorm:"instruct_flag" json:"instruct_flag"`         // 指令下发方式 1:单条下发 2:多条下发
	SimCode          string `gorm:"sim_code" json:"sim_code"`                   // SIM卡ICCID
	BindStatus       int    `gorm:"bind_status" json:"bind_status"`             // 绑定状态 0:未绑定 1:已绑定
	LineStatus       int    `gorm:"line_status" json:"line_status"`             // 设备状态: 0:离线 1:在线
	DeviceAddress    string `gorm:"device_address" json:"device_address"`       // 设备地址
	DeviceCoordinate string `gorm:"device_coordinate" json:"device_coordinate"` // 设备坐标信息
	BindGroup        int    `gorm:"bind_group" json:"bind_group"`               // 绑定分组id

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
	FilterValue   float32 `json:"filterValue"`   //过滤值
	LeftValueType int     `json:"leftValueType"` //左侧值
	Level         int     `json:"level"`         //报警等级
	Operator      string  `json:"operator"`      //比对条件
}

type ShakeLimit struct { //防抖
	Enabled   int    `json:"enabled"`   //是否开启防抖
	Time      int    `json:"time"`      //时间间隔(秒)
	Threshold int    `json:"threshold"` //触发阈值(次)
	Handle    string `json:"Handle"`    //是否第一次满足条件就触发

}

// Device 结构体
type Device struct {
	GuId         string `json:"guid"`
	Name         string `json:"name"`
	ProductType  int    `json:"productType"`
	ProductModel string `json:"productModel"`

	FCode *Ft      `json:"fCode"`
	Salve []*Slave `json:"salve"`
}

func (device *Device) IsType() bool {

	if device.ProductType == 1 { //工业
		return true
	} else {
		return false
	}
}

// Detector 结构体
// 如果含有time.Time 请自行import time包
type Detector struct {
	Id                string     `json:"id" gorm:"primarykey；column:id;"`
	DeviceId          string     `json:"deviceId" form:"deviceId" gorm:"column:device_id;comment:设备表id;size:100;"`
	Name              string     `json:"name" form:"name" gorm:"column:name;comment:从机设备名称;size:100;"`
	OtherName         string     `json:"otherName" form:"otherName" gorm:"column:other_name;comment:从机设备别名;size:100;"`
	SlaveAddress      int        `json:"slaveAddress" form:"slaveAddress" gorm:"column:slave_address;comment:modbus从机地址;size:5;"`
	Status            int        `json:"status" form:"status" gorm:"column:status;comment:状态（0：离线；1：在线）;size:10;"`
	AttributeId       string     `json:"attributeId" form:"attributeId" gorm:"column:attribute_id;comment:属性表的id;size:100;"`
	SlaveDeviceSwitch int        `json:"slaveDeviceSwitch" form:"slaveDeviceSwitch" gorm:"column:slave_device_switch;comment:从机设备开关（1--开，0--关）;size:10;"`
	Remarks           string     `json:"remarks" form:"remarks" gorm:"column:remarks;comment:从机有关描述;size:255;"`
	CreateTime        *time.Time `json:"createTime" form:"createTime" gorm:"column:create_time;comment:创建时间;"`
	UpdateTime        *time.Time `json:"updateTime" form:"updateTime" gorm:"column:update_time;comment:修改时间;"`
	CreateBy          string     `json:"createBy" form:"createBy" gorm:"column:create_by;comment:创建者;size:64;"`
	UpdateBy          string     `json:"updateBy" form:"updateBy" gorm:"column:update_by;comment:更新人;size:64;"`
	//Attribute         Attribute  `json:"attribute" gorm:"foreignKey:AttributeId;references:AttributeId;comment:属性数据"`
}

// TableName Detector 表名
func (Detector) TableName() string {
	return "detector"
}

// Attribute 结构体
// 如果含有time.Time 请自行import time包
type Attribute struct {
	Id           string     `json:"id" gorm:"primarykey；column:id;"`
	DeviceId     string     `json:"deviceId" form:"deviceId" gorm:"column:device_id;comment:设备id;size:100;"`
	ProductId    string     `json:"productId" form:"productId" gorm:"column:product_id;comment:产品id;size:100;"`
	Name         string     `json:"name" form:"name" gorm:"column:name;comment:属性名称;size:100;"`
	Code         string     `json:"code" form:"code" gorm:"column:code;comment:标识符;size:100;"`
	SourceType   int        `json:"sourceType" form:"sourceType" gorm:"column:source_type;comment:来源类型（1：主动上报；2：服务器采集）;size:10;"`
	Frequency    int        `json:"frequency" form:"frequency" gorm:"column:frequency;comment:采集频率，只有来源类型为服务器采集，才有采集频率（1：30秒；2：1分钟；3：5分钟）;size:10;"`
	DataType     int        `json:"dataType" form:"dataType" gorm:"column:data_type;comment:数据类型（1：整数；2：字符；3：文本；4：浮点）;size:10;"`
	UnitTypeId   int        `json:"unitTypeId" form:"unitTypeId" gorm:"column:unit_type_id;comment:单位类型id;size:10;"`
	Digit        int        `json:"digit" form:"digit" gorm:"column:digit;comment:小数位数（0：没有小数位数；1:1个小数位数；2:2个小数位数；3:3个小数位数；4:4个小数位数）;size:10;"`
	FunctionCode int        `json:"functionCode" form:"functionCode" gorm:"column:function_code;comment:功能码（1：01H；2：02H；3：03H；4：04H；5：05H；6：06H；7：0FH；8：10H）;size:10;"`
	RegAddress   int        `json:"regAddress" form:"regAddress" gorm:"column:reg_address;comment:寄存器地址（地址偏移）;size:10;"`
	IsRead       int        `json:"isRead" form:"isRead" gorm:"column:is_read;comment:是否只读（0:否；1：是）;size:10;"`
	ErrorMessage string     `json:"errorMessage" form:"errorMessage" gorm:"column:error_message;comment:报错信息;size:255;"`
	Remarks      string     `json:"remarks" form:"remarks" gorm:"column:remarks;comment:描述;size:255;"`
	CreateTime   *time.Time `json:"createTime" form:"createTime" gorm:"column:create_time;comment:创建时间;"`
	UpdateTime   *time.Time `json:"updateTime" form:"updateTime" gorm:"column:update_time;comment:修改时间;"`
	CreateBy     string     `json:"createBy" form:"createBy" gorm:"column:create_by;comment:创建者;size:64;"`
	UpdateBy     string     `json:"updateBy" form:"updateBy" gorm:"column:update_by;comment:更新人;size:64;"`
}

// TableName Attribute 表名
func (Attribute) TableName() string {
	return "attribute"
}

// Condition 结构体
// 如果含有time.Time 请自行import time包
type Condition struct {
	Id            string     `json:"id" gorm:"primarykey；column:id;"`
	AlarmRuleId   string     `json:"alarmRuleId" form:"alarmRuleId" gorm:"column:alarm_rule_id;comment:告警规则id;size:100;"`
	ConditionType string     `json:"conditionType" form:"conditionType" gorm:"column:condition_type;comment:条件类型（1：属性）;size:10;"`
	AttributeId   string     `json:"attributeId" form:"attributeId" gorm:"column:attribute_id;comment:属性的id;size:100;"`
	Type          int        `json:"type" form:"type" gorm:"column:type;comment:类型（1：最新值；2：平均值；3：最大值；4：最小值；5：和值；6：值有变化；7：无值）;size:10;"`
	SymbolType    string     `json:"symbolType" form:"symbolType" gorm:"column:symbol_type;comment:符号类型（1：=；2：>；3：<；4：<>；5：>=；6：<=）;size:10;"`
	Values        string     `json:"values" form:"values" gorm:"column:values;comment:作对比用的值;size:30;"`
	CreateTime    *time.Time `json:"createTime" form:"createTime" gorm:"column:create_time;comment:创建时间;"`
	UpdateTime    *time.Time `json:"updateTime" form:"updateTime" gorm:"column:update_time;comment:修改时间;"`
	CreateBy      string     `json:"createBy" form:"createBy" gorm:"column:create_by;comment:创建者;size:64;"`
	UpdateBy      string     `json:"updateBy" form:"updateBy" gorm:"column:update_by;comment:更新人;size:64;"`
}

// TableName Condition 表名
func (Condition) TableName() string {
	return "condition"
}

// AlarmRule 结构体
// 如果含有time.Time 请自行import time包
type AlarmRule struct {
	Id           string     `json:"id" gorm:"primarykey；column:id;"`
	DeviceId     string     `json:"deviceId" form:"deviceId" gorm:"column:device_id;comment:设备表id;size:100;"`
	Name         string     `json:"name" form:"name" gorm:"column:name;comment:告警规则名称;size:100;"`
	AlarmLevel   int        `json:"alarmLevel" form:"alarmLevel" gorm:"column:alarm_level;comment:告警级别（1--低，2--中，3--高）;size:10;"`
	ProductId    string     `json:"productId" form:"productId" gorm:"column:product_id;comment:产品id;size:100;"`
	EnableStatus int        `json:"enableStatus" form:"enableStatus" gorm:"column:enable_status;comment:启用状态（0--未启用；1--已启用）;size:10;"`
	Remarks      string     `json:"remarks" form:"remarks" gorm:"column:remarks;comment:告警规则有关描述;size:100;"`
	Condition    int        `json:"condition" form:"condition" gorm:"column:condition;comment:满足哪种条件，触发告警（1：任意；2：所有）;size:10;"`
	IsShake      int        `json:"isShake" form:"isShake" gorm:"column:is_shake;comment:是否开启防抖（0：否；1：是）;size:10;"`
	Within       int        `json:"within" form:"within" gorm:"column:within;comment:在···时间里（开启防抖才有值）秒;size:10;"`
	Num          int        `json:"num" form:"num" gorm:"column:num;comment:发生次数（开启防抖才有值）;size:10;"`
	IsFirst      int        `json:"isFirst" form:"isFirst" gorm:"column:is_first;comment:处理哪一次（1：第一次；2：最后一次）（开启防抖才有值）;size:10;"`
	CreateTime   *time.Time `json:"createTime" form:"createTime" gorm:"column:create_time;comment:创建时间;"`
	UpdateTime   *time.Time `json:"updateTime" form:"updateTime" gorm:"column:update_time;comment:修改时间;"`
	CreateBy     string     `json:"createBy" form:"createBy" gorm:"column:create_by;comment:创建者;size:64;"`
	UpdateBy     string     `json:"updateBy" form:"updateBy" gorm:"column:update_by;comment:更新人;size:64;"`
}

// TableName AlarmRule 表名
func (AlarmRule) TableName() string {
	return "alarm_rule"
}

// ExecuteAction 结构体
// 如果含有time.Time 请自行import time包
type ExecuteAction struct {
	Id               string     `json:"id" gorm:"primarykey；column:id;"`
	Type             string     `json:"type" form:"type" gorm:"column:type;comment:执行动作类型（1：消息通知；2：设备通知）;size:10;"`
	AlarmRuleId      string     `json:"alarmRuleId" form:"alarmRuleId" gorm:"column:alarm_rule_id;comment:关联告警规则id;size:30;"`
	NotifyType       string     `json:"notifyType" form:"notifyType" gorm:"column:notify_type;comment:通知类型(1:语音；2：短信；3：微信等);size:10;"`
	NotifyConfigId   string     `json:"notifyConfigId" form:"notifyConfigId" gorm:"column:notify_config_id;comment:关联通知配置表id;size:30;"`
	NotifyTemplateId string     `json:"notifyTemplateId" form:"notifyTemplateId" gorm:"column:notify_template_id;comment:关联通知模板id;size:30;"`
	CreateTime       *time.Time `json:"createTime" form:"createTime" gorm:"column:create_time;comment:创建时间;"`
	UpdateTime       *time.Time `json:"updateTime" form:"updateTime" gorm:"column:update_time;comment:修改时间;"`
	CreateBy         string     `json:"createBy" form:"createBy" gorm:"column:create_by;comment:创建者;size:64;"`
	UpdateBy         string     `json:"updateBy" form:"updateBy" gorm:"column:update_by;comment:更新人;size:64;"`
}

// TableName ExecuteAction 表名
func (ExecuteAction) TableName() string {
	return "execute_action"
}

type NotifyConfig struct {
	Id            string `json:"id" gorm:"primarykey；column:id;"`
	Name          string `json:"name" form:"name" gorm:"column:name;comment:告警规则名称;size:100;"`
	Type          string `json:"type" form:"type" gorm:"column:type;comment:执行动作类型（1：消息通知；2：设备通知）;size:10;"`
	Provider      string `json:"provider" form:"provider" gorm:"column:provider;comment:服务商(aliyun:阿里云；aliyunSms：阿里云短信服务；wechat：微信通知);size:10;"`
	Configuration string `json:"configuration" form:"configuration" gorm:"column:configuration;comment:配置内容;size:10;"`
}

// TableName NotifyConfig 表名
func (NotifyConfig) TableName() string {
	return "notify_config"
}

type NotifyTemplate struct {
	Id             string `json:"id" gorm:"primarykey；column:id;"`
	NotifyConfigId string `json:"notifyConfigId" gorm:"primarykey；column:notify_config_id;"`
	Name           string `json:"name" form:"name" gorm:"column:name;comment:模板名称;size:100;"`
	Type           string `json:"type" form:"type" gorm:"column:type;comment:执行动作类型（1：消息通知；2：设备通知）;size:10;"`
	Provider       string `json:"provider" form:"provider" gorm:"column:provider;comment:服务商(aliyun:阿里云；aliyunSms：阿里云短信服务；wechat：微信通知);size:10;"`
	Template       string `json:"template" form:"template" gorm:"column:template;comment:模板内容;size:10;"`
}

// TableName NotifyConfig 表名
func (NotifyTemplate) TableName() string {
	return "notify_template"
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
	SlaveId    byte   `json:"slaveId"`
	SlaveName  string `json:"slaveName"`
	Alarm      *Alarm `json:"alarm"`
	DataTime   time.Time
	LineStatus string
}

//type Trigger struct { //触发条件
//	Type     string `json:"type"`     //触发条件类型
//	ModelId  string `json:"modelId"`  //属性ID
//	Operator string `json:"operator"` //条件
//	Val      uint16 `json:"val"`      //数据值
//}

type Action struct { //执行动作
	Type       string `json:"type"`       //执行动作类型
	NotifyType string `json:"notifyType"` //通知类型
	NotifierId string `json:"notifierId"` //通知配置ID
	TemplateId string `json:"templateId"` //通知模版ID
}

//type ShakeLimit struct { //防抖
//	Enabled    bool `json:"enabled"`    //是否开启防抖
//	Time       int  `json:"time"`       //时间间隔(秒)
//	Threshold  int  `json:"threshold"`  //触发阈值(次)
//	AlarmFirst bool `json:"alarmFirst"` //是否第一次满足条件就触发
//}

type ListenMsg struct {
	ListenType int
	RemoteAddr string
	Guid       string
	Command    Comm
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

type Interface interface {
	AlarmRule(slaveId byte, data float32, fcode uint8, info *Device)
	Trigger(data float32)
	Action(guid, status string, level int, data float32, slaveId byte)
}

type Alarm struct {
	Triggers   []*Trigger  `json:"trigger"`    //条件
	ShakeLimit *ShakeLimit `json:"shakeLimit"` //防抖动配置
}

func (engine *Alarm) AlarmRule(slaveId byte, data float32, fcode uint8, info *Device) {
	if info.IsType() {
		switch data { //是否工业产品
		// 10000     探测器内部错误
		case consts.InternalError:
			engine.Action(info.GuId, consts.ALARM, consts.Internal, data, slaveId)
			break

		// 20000     通讯错误
		case consts.CommunicationError:
			engine.Action(info.GuId, consts.ALARM, consts.Communication, data, slaveId)
			break

		// 30000     主机未连接探测器、主机屏蔽探测器
		case consts.ShieldError:
			engine.Action(info.GuId, consts.ALARM, consts.Shield, data, slaveId)
			break

			// 65535     探头故障
		case consts.SlaveHitchError:
			engine.Action(info.GuId, consts.ALARM, consts.SlaveHitch, data, slaveId)
			break
		}
		engine.Trigger(slaveId, data, info)
	} else {
		switch fcode {
		case consts.ReadCode:
			engine.Trigger(slaveId, data, info)
			break
		case consts.HomeHitchError:
			//故障（若为0，是传感器低故障报警，若为1，是传感器高故障报警，若为2，是传感器寿命报警）
			if data == 0 {
				engine.Action(info.GuId, consts.ALARM, consts.LowHitch, data, slaveId)
			} else if data == 1 {
				engine.Action(info.GuId, consts.ALARM, consts.HighHitch, data, slaveId)
			} else if data == 2 {
				engine.Action(info.GuId, consts.ALARM, consts.Life, data, slaveId)
			}
			break
		case consts.HomeHighError:
			engine.Action(info.GuId, consts.ALARM, consts.High, data, slaveId)
			break
		case consts.HomeLowError:
			engine.Action(info.GuId, consts.ALARM, consts.Low, data, slaveId)
			break
		}

	}
}

func (engine *Alarm) Trigger(slaveId byte, data float32, info *Device) {

	for _, trigger := range engine.Triggers { //循环告警触发条件

		switch trigger.Operator { //TODO 判断比对条件(任意) 触发条件满足条件中任意一个即可触发  高报优先
		case consts.EQ: //=
			if data == trigger.FilterValue {
				engine.Action(info.GuId, consts.ALARM, trigger.Level, data, slaveId)
				return
			}
		case consts.NOT:
			if data != trigger.FilterValue {
				engine.Action(info.GuId, consts.ALARM, trigger.Level, data, slaveId)
				return
			}
		case consts.GT:
			if data > trigger.FilterValue {
				engine.Action(info.GuId, consts.ALARM, trigger.Level, data, slaveId)
				return
			}
		case consts.LT:
			if data < trigger.FilterValue {
				engine.Action(info.GuId, consts.ALARM, trigger.Level, data, slaveId)
				return
			}
		case consts.GTE:
			if data >= trigger.FilterValue {
				engine.Action(info.GuId, consts.ALARM, trigger.Level, data, slaveId)
				return
			}
		case consts.LTE:
			if data <= trigger.FilterValue {
				engine.Action(info.GuId, consts.ALARM, trigger.Level, data, slaveId)
				return
			}
		}
	}
	engine.Action(info.GuId, consts.DATA, consts.Normal, data, slaveId)
}

func (engine *Alarm) Action(guid, status string, level int, data float32, slaveId byte) {

	device.DataChan <- &device.DeviceMsg{
		Ts:       time.Now(),
		DataType: consts.DATA,
		Level:    level,
		DeviceId: guid,
		SlaveId:  int(slaveId),
		Data:     data,
	}
	if status == consts.ALARM {
		device.AlarmChan <- &device.DeviceMsg{
			Ts:       time.Now(),
			DataType: consts.ALARM,
			Level:    level,
			DeviceId: guid,
			SlaveId:  int(slaveId),
			Data:     data,
		}
	}

}
