package db

import (
	"database/sql"
	"gorm.io/gorm"
)

var (
	DB *gorm.DB
	Td *sql.DB
)

func InitGorm(db *gorm.DB) {
	DB = db
}

func InitTdengine(db *sql.DB) {
	Td = db
}
