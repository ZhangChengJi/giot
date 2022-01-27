package model

import (
	"github.com/panjf2000/gnet"
	"time"
)

type RemoteData struct {
	Frame      []byte
	RemoteAddr string
}
type RegisterData struct {
	C gnet.Conn
	D []byte
}

// Device 结构体
type Device struct {
	Id                string     `json:"id" gorm:"primarykey；column:id;"`
	Name              string     `json:"name" form:"name" gorm:"column:name;comment:设备名称;size:100;"`
	ProductId         string     `json:"productId" form:"productId" gorm:"column:product_id;comment:产品ID;size:100;"`
	Type              int        `json:"type" form:"type" gorm:"column:type;comment:设备类型（1--直连设备，2--网关设备，3--网关子设备）;size:10;"`
	EnableStatus      int        `json:"enableStatus" form:"enableStatus" gorm:"column:enable_status;comment:启用状态（0：未启用；1：已启用）;size:10;"`
	OnlineStatus      int        `json:"onlineStatus" form:"onlineStatus" gorm:"column:online_status;comment:在线状态（0--离线，1--注册，2--上数和报警）;size:10;"`
	ProCode           string     `json:"proCode" form:"proCode" gorm:"column:pro_code;comment:产品编码;size:100;"`
	DeviceGroupId     string     `json:"deviceGroupId" form:"deviceGroupId" gorm:"column:device_group_id;comment:设备组id;size:100;"`
	DeptId            int        `json:"deptId" form:"deptId" gorm:"column:dept_id;comment:机构id;size:10;"`
	Remarks           string     `json:"remarks" form:"remarks" gorm:"column:remarks;comment:设备有关描述;size:100;"`
	LongitudeLatitude string     `json:"longitudeLatitude" form:"longitudeLatitude" gorm:"column:longitude_latitude;comment:经度和纬度;size:100;"`
	Address           string     `json:"address" form:"address" gorm:"column:address;comment:设备详细地址;size:100;"`
	CreateTime        *time.Time `json:"createTime" form:"createTime" gorm:"column:create_time;comment:创建时间;"`
	UpdateTime        *time.Time `json:"updateTime" form:"updateTime" gorm:"column:update_time;comment:修改时间;"`
	CreateBy          string     `json:"createBy" form:"createBy" gorm:"column:create_by;comment:创建者;size:64;"`
	UpdateBy          string     `json:"updateBy" form:"updateBy" gorm:"column:update_by;comment:更新人;size:64;"`
}

func (Device) TableName() string {
	return "transfer"
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

//type Device struct {
//	Guid            string     //设备ID
//	Name            string     //设备名称
//	EnableStatus    int        //设备启用状态
//	ProCode         string     //产品编码
//	MessageProtocol int        //消息协议
//	DetectorList    []Detector //从机
//	FunctionCode    []byte     //下发指令
//	Rule            []Rule
//}
//type Detector struct {
//	Name         string
//	SlaveAddress string
//	RegAddress   string
//}

type TimerActive struct {
	Guid string `json:"guId"`
	Ft   []*Ft  `json:"ft"`
}
type Ft struct {
	Tm    time.Duration `json:"tm"`
	FCode [][]byte      `json:"fCode"`
}

type Slave struct {
	ProductId   string `json:"productId"`   //产品ID
	ProductName string `json:"productName"` //产品名称
	DeviceId    string `json:"deviceId"`    //设备Id
	DeviceName  string `json:"deviceName"`  //设备名称
	SlaveId     byte   `json:"slaveId"`
	SlaveName   string `json:"slaveName"`
	AttributeId string `json:"attributeId"`
}
type Alarm struct {
	Name        string      `json:"name"`        //告警名称
	ProductId   string      `json:"productId"`   //产品ID
	ProductName string      `json:"productName"` //产品名称
	DeviceId    string      `json:"deviceId"`    //设备Id
	DeviceName  string      `json:"deviceName"`  //设备名称
	ShakeLimit  *ShakeLimit `json:"shakeLimit"`  //防抖动配置
	ModelId     []string    `json:"ModelId"`     //属性ID
	Triggers    []*Trigger  `json:"triggers"`    //触发条件
	Actions     []*Action   `json:"actions"`     //执行动作
}

type Trigger struct { //触发条件
	Type     string `json:"type"`     //触发条件类型
	ModelId  string `json:"modelId"`  //属性ID
	Operator string `json:"operator"` //条件
	Val      []byte `json:"val"`      //数据值
}

type Action struct { //执行动作
	Type       string `json:"type"`       //执行动作类型
	NotifyType string `json:"notifyType"` //通知类型
	NotifierId string `json:"notifierId"` //通知配置ID
	TemplateId string `json:"templateId"` //通知模版ID
}

type ShakeLimit struct { //防抖
	Enabled    bool `json:"enabled"`    //是否开启防抖
	Time       int  `json:"time"`       //时间间隔(秒)
	Threshold  int  `json:"threshold"`  //触发阈值(次)
	AlarmFirst bool `json:"alarmFirst"` //是否第一次满足条件就触发
}

type ListenMsg struct {
	ListenType int
	RemoteAddr string
	Guid       string
	Command    Comm
}

type Comm int8
