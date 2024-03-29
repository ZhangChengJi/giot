module giot

go 1.13

replace google.golang.org/grpc => google.golang.org/grpc v1.26.0

replace github.com/coreos/bbolt => go.etcd.io/bbolt v1.3.5

require (
	github.com/RussellLuo/timingwheel v0.0.0-20201029015908-64de9d088c74
	github.com/allegro/bigcache/v3 v3.0.2
	github.com/coreos/bbolt v0.0.0-00010101000000-000000000000 // indirect
	github.com/coreos/etcd v3.3.25+incompatible // indirect
	github.com/coreos/go-systemd v0.0.0-20191104093116-d3cd4ed1dbcf // indirect
	github.com/coreos/pkg v0.0.0-20180928190104-399ea9e2e55f // indirect
	github.com/dgrijalva/jwt-go v3.2.0+incompatible // indirect
	github.com/dustin/go-humanize v1.0.0 // indirect
	github.com/eclipse/paho.mqtt.golang v1.3.5
	github.com/go-redis/redis/v8 v8.11.5
	github.com/google/uuid v1.3.0 // indirect
	github.com/grpc-ecosystem/go-grpc-middleware v1.3.0 // indirect
	github.com/grpc-ecosystem/go-grpc-prometheus v1.2.0 // indirect
	github.com/jonboulle/clockwork v0.2.2 // indirect
	github.com/json-iterator/go v1.1.12
	github.com/kr/text v0.2.0 // indirect
	github.com/lestrrat-go/file-rotatelogs v2.4.0+incompatible
	github.com/lestrrat-go/strftime v1.0.6 // indirect
	github.com/lithammer/dedent v1.1.0
	github.com/niemeyer/pretty v0.0.0-20200227124842-a10e7caefd8e // indirect
	github.com/onsi/gomega v1.20.2 // indirect
	github.com/orcaman/concurrent-map v1.0.0
	github.com/panjf2000/gnet v1.6.6
	github.com/panjf2000/gnet/v2 v2.1.1
	github.com/shiningrush/droplet v0.3.0
	github.com/shopspring/decimal v1.3.1
	github.com/soheilhy/cmux v0.1.5 // indirect
	github.com/spf13/cobra v1.3.0
	github.com/spf13/viper v1.10.1
	github.com/stretchr/testify v1.7.0
	github.com/taosdata/driver-go/v2 v2.0.1
	github.com/tmc/grpc-websocket-proxy v0.0.0-20220101234140-673ab2c3ae75 // indirect
	github.com/vmihailenco/msgpack/v5 v5.3.5
	github.com/xiang90/probing v0.0.0-20190116061207-43a291ad63a2 // indirect
	go.etcd.io/bbolt v1.3.6 // indirect
	go.etcd.io/etcd v3.3.25+incompatible
	go.uber.org/zap v1.21.0
	golang.org/x/crypto v0.0.0-20220131195533-30dcbda58838 // indirect
	golang.org/x/text v0.3.7
	gopkg.in/check.v1 v1.0.0-20200227125254-8fa46927fb4f // indirect
	gorm.io/driver/mysql v1.2.3
	gorm.io/gorm v1.22.4
	sigs.k8s.io/yaml v1.3.0 // indirect
)
