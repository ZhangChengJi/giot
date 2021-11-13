package conf

const (
	EnvPROD  = "prod"
	EnvBETA  = "beta"
	EnvDEV   = "dev"
	EnvLOCAL = "local"
	EnvTEST  = "test"

	WebDir = "html/"
)

var (
	ENV string

	ServerHost    = "127.0.0.1"
	ServerPort    = 80
	ErrorLogLevel = "warn"
	ErrorLogPath  = "logs/error.log"
	AccessLogPath = "logs/access.log"
)
