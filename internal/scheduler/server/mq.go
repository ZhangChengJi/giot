package server

import (
	"giot/conf"
	broker "giot/internal/scheduler/mqtt"
	"giot/pkg/log"
	"giot/pkg/mqtt"
)

func (s *server) setupMqtt() error {
	conf.MqttConfig.ClientId = "scheduler"
	if client, err := mqtt.New(conf.MqttConfig); err != nil {
		log.Errorf("init mqtt client fail: %w", err)
		return err
	} else {
		broker.InitMqtt(client)
	}

	return nil
}
