package server

import (
	"giot/conf"
	"giot/pkg/etcd"
	"giot/pkg/log"
)

func (s *server) setupEtcd() error {
	if err := etcd.InitETCDClient(conf.ETCDConfig); err != nil {
		log.Errorf("init etcd client fail: %w", err)
		return err
	}
	return nil
}
