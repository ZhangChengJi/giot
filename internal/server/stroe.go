package server

import (
	"giot/internal/conf"
	"giot/internal/log"
	"giot/internal/storage"
)

func (s *server) setupStore() error {
	if err := storage.InitETCDClient(conf.ETCDConfig); err != nil {
		log.Errorf("init etcd client fail: %w", err)
		return err
	}
	return nil
}
