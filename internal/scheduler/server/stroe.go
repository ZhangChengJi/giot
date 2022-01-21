package server

import (
	"giot/internal/scheduler/conf"
	"giot/internal/scheduler/log"
	"giot/internal/scheduler/storage"
)

func (s *server) setupStore() error {
	if err := storage.InitETCDClient(conf.ETCDConfig); err != nil {
		log.Errorf("init etcd client fail: %w", err)
		return err
	}
	return nil
}
