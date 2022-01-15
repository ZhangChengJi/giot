package db

import (
	"gorm.io/gorm"
)

var (
	DB *gorm.DB
)

func InitGorm(db *gorm.DB) {
	DB = db
}
