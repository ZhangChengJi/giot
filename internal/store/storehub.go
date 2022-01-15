package store

import (
	"fmt"
	"giot/internal/core/model"
	"giot/internal/engine"
	"giot/internal/storage"
)

type HubKey string

const (
	HubKeyDevicePrefix HubKey = "logic"
)

type DeviceStore struct {
	Stg storage.Interface
	eng engine.Interface
}

var (
	deviceHub = map[HubKey]*DeviceStore{} //数据处理接口集合
)

func IsDevice(key HubKey) bool {
	_, ok := deviceHub[key]
	return ok
}

func GetDevice(key HubKey) *DeviceStore {
	if s, ok := deviceHub[key]; ok {
		return s
	}
	panic(fmt.Sprintf("no logic with key: %s", key))
}

func WithPrefix(str string) HubKey {
	return HubKeyDevicePrefix + "/" + HubKey(str)
}
func NewDevice(key HubKey, device model.Device) error {
	//deviceHub[key] = logic
	return nil
}
