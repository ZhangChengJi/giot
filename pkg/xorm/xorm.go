package xorm

import (
	"fmt"
	"giot/internal/conf"
	_ "github.com/lib/pq"
	"github.com/xormplus/xorm"
)

func New(conf *conf.Postgres) (*xorm.Engine, error) {
	dataSource := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable", conf.Username, conf.Password, conf.Host, conf.Port, conf.Database)
	engine, err := xorm.NewPostgreSQL(dataSource)
	if err != nil {
		return nil, err
	}
	engine.ShowSQL(conf.ShowSql) //是否开启打印sql
	engine.SetMaxIdleConns(conf.MaxIdleConns)
	engine.SetMaxOpenConns(conf.MaxOpenConns)
	err = engine.Ping()
	if err != nil {
		return nil, err
	}
	return engine, nil
}
