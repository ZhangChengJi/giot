package server

import (
	"giot/internal/conf"
	"giot/internal/db"
	"giot/internal/log"
	"giot/pkg/gorm"
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
