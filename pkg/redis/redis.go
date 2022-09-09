package redis

import (
	"giot/conf"
	"giot/pkg/log"
	"github.com/go-redis/redis"
)

func New(conf *conf.Redis) (*redis.Client, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     conf.Addr,
		Password: conf.Password,
		DB:       conf.Db,
	})
	ping, err := client.Ping().Result()
	if err != nil {
		log.Sugar.Errorf("redis连接错误%v", err.Error())
		return nil, err
	}
	log.Sugar.Infof("redis连接成功%v", ping)
	return client, nil
}
