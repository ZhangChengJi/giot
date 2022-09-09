package consts

var (
	ActionAll   = "all"
	ActionCode  = "code"
	ActionSlave = "slave"
	ActionAlarm = "alarm"
)

const (
	Update = "update"
	Add    = "add"
	Delete = "delete"
)

const (
	Normal        = iota //正常
	High                 //高报
	Low                  //低报
	Internal             //探测器内部错误
	Communication        //通讯错误
	Shield               //主机屏蔽探测器
	SlaveHitch           //探头故障
	LowHitch             //家报传感器低故障
	HighHitch            //家报传感器高故障
	Life                 //家报寿命
	DATA          = "data"
	ALARM         = "alarm"
	HITCH         = "hitch"
	SMS           = "sms"
	VOICE         = "voice"
	WECHAT        = "wechat"
	ONLINE        = "online"  //上线
	OFFLINE       = "offline" //下线

	EQ    = "eq"  //==
	NOT   = "not" //<>
	GT    = "gt"  //>
	LT    = "lt"  //<
	GTE   = "gte" //>=
	LTE   = "lte" //<=
	FIRST = "first"
	LAST  = "last"

	//+++++++++++++++工业+++++++++++++++++++
	InternalError      = 10000 //探测器内部错误
	CommunicationError = 20000 //通讯错误
	ShieldError        = 30000 //主机未连接探测器、主机屏蔽探测器
	SlaveHitchError    = 65535 //探头故障

	//+++++++++++++家用++++++++++++++++++++
	HomeHitchError = 63 //家用故障功能码 //故障（数据值若为0，是传感器低故障报警，若为1，是传感器高故障报警，若为2，是传感器寿命报警）
	HomeHighError  = 65 //家用高报警功能码
	HomeLowError   = 66 //家用低报警功能码
	ReadCode       = 03
)
