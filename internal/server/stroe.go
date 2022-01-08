package server

//import (
//	"giot/internal/conf"
//	"giot/internal/log"
//	"giot/internal/stroage"
//	"giot/pkg/xorm"
//)
//
//func (s *server) setupDb() error {
//	if engine, err := xorm.New(conf.PostgresConfig); err != nil {
//		log.Errorf("postgres connection fail: %s", err.Error())
//		return err
//	} else {
//		err := stroage.InitXorm(engine)
//		if err != nil {
//			return err
//		}
//	}
//	return nil
//
//}

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
