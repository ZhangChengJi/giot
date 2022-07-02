package server

import (
	"giot/conf"
	broker "giot/internal/virtual/mqtt"
	"giot/pkg/log"
	"giot/pkg/mqtt"
	"math/rand"
	"strconv"
	"time"
)

func (s *server) setupMqtt() error {
	conf.MqttConfig.ClientId = "virtual" + strconv.Itoa(rand.New(rand.NewSource(time.Now().UnixNano())).Int())
	if client, err := mqtt.New(conf.MqttConfig); err != nil {
		log.Sugar.Errorf("init mqtt client fail: %w", err)
		return err
	} else {
		broker.InitMqtt(client)
	}

	return nil
}
