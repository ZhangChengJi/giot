package server

import (
	"giot/internal/scheduler/conf"
	"giot/internal/scheduler/db"
	"giot/internal/scheduler/log"
	"giot/pkg/gorm"
	"giot/pkg/tdengine"
)

func (s *server) setupDB() error {
	if engine, err := gorm.New(conf.MysqlConfig); err != nil {
		log.Errorf("mysql connection fail: %s", err.Error())
		return err
	} else {
		db.InitGorm(engine)
	}
	return nil

}

func (s *server) setupTdengine() error {
	if engine, err := tdengine.New(conf.TdengineConfig); err != nil {
		log.Errorf("tdengine connection fail: %s", err.Error())
		return err
	} else {
		db.InitTdengine(engine)
	}
	return nil
}
