package server

import (
	"giot/conf"
	"giot/pkg/etcd"
	"giot/pkg/log"
)

func (s *server) setupStore() error {
	if err := etcd.InitETCDClient(conf.ETCDConfig); err != nil {
		log.Sugar.Errorf("init etcd client fail: %w", err)
		return err
	}

	return nil
}
