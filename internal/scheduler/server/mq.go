package server

import (
	"giot/internal/scheduler/log"
	broker "giot/internal/scheduler/mqtt"
	"giot/pkg/mqtt"
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
