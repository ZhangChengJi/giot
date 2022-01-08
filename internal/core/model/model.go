package model

type Device struct {
	Guid            string     //设备ID
	Name            string     //设备名称
	EnableStatus    int        //设备启用状态
	ProCode         string     //产品编码
	MessageProtocol int        //消息协议
	DetectorList    []Detector //从机
	FunctionCode    []byte     //下发指令
	Rule            []Rule
}
type Detector struct {
	Name         string
	SlaveAddress string
	RegAddress   string
}
type Rule struct {
	ShakeLimit ShakeLimit
	Triggers   []trigger //触发条件
	Action     []Action  //执行动作
}
type trigger struct { //触发条件
	ModelId  string //属性ID
	Operator string //条件

}

func Operator(vale string) {
	const (
		eq   = "="
		not  = "!="
		gt   = ">"
		lt   = "<"
		gte  = ">="
		lte  = "<="
		like = "like"
	)

}

type Action struct { //执行动作
}

type ShakeLimit struct { //防抖
	enabled    bool //是否开启防抖
	time       int  //时间间隔(秒)
	threshold  int  //触发阈值(次)
	alarmFirst bool //是否第一次满足条件就触发
}

type DtuMsg struct {
	host string
}
