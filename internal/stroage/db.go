package stroage

import (
	"giot/internal/log"
	"github.com/xormplus/xorm"
)

var (
	DB *xorm.Engine
)

func InitXorm(db *xorm.Engine) error {
	err := db.RegisterSqlMap(xorm.Xml("/Users/zhangchengji/GoProjects/giot/internal/stroage/sql/postgres", ".xml"))
	if err != nil {
		log.Errorf("Register sqlmap fail:%s", err.Error())
		return err
	}
	DB = db
	return nil
}
