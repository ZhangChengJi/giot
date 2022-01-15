package gorm

import (
	"fmt"
	"giot/internal/conf"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func New(conf *conf.Mysql) (*gorm.DB, error) {
	dataSource := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local", conf.Username, conf.Password, conf.Host, conf.Port, conf.Database)
	db, err := gorm.Open(mysql.Open(dataSource), &gorm.Config{})
	if err != nil {
		return nil, err

	}
	sqlDB, _ := db.DB()
	sqlDB.SetMaxIdleConns(conf.MaxIdleConns)
	sqlDB.SetMaxOpenConns(conf.MaxOpenConns)
	return db, nil
}
