package logic

import (
	"giot/internal/model"
	"giot/internal/scheduler/db"
	"giot/pkg/log"
)

//设备↕上下线

func Online(guid string) {
	var device model.Device
	err := db.DB.First(guid, device).Error
	if err != nil {
		log.Errorf("online guid:%s not found", guid)
		return
	}
	err = db.DB.Model(&device).Update("online_status", 1).Error
	if err != nil {
		log.Errorf("online guid:%s update failed", guid)
		return
	}
}
func Offline(guid string) {
	var device model.Device
	err := db.DB.First(guid, device).Error
	if err != nil {
		log.Errorf("online guid:%s not found", guid)
		return
	}
	err = db.DB.Model(&device).Update("online_status", 0).Error
	if err != nil {
		log.Errorf("online guid:%s update failed", guid)
		return
	}
}
