package server

import (
	"giot/pkg/log"
	mqtt "github.com/eclipse/paho.mqtt.golang"
)

func (s *server) setupMqtt() error {

	if client, err := mqtt.NewClient(); err != nil {
		log.Errorf("init mqtt client fail: %w", err)
		return err
	} else {
		broker.InitMqtt(client)
	}

	return nil
}
