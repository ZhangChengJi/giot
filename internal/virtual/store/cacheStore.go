package store

import "giot/pkg/log"

type Interface interface {
}

type CacheStore struct {
	Device DeviceStoreIn
	Guid   GuidStoreIn
	Slave  SlaveStoreIn
	Timer  DeviceTimerIn
	Line   LineStoreIn
}

func New() *CacheStore {
	store := &CacheStore{
		Device: NewDeviceStore(),
		Guid:   NewGuidStore(),
		Slave:  NewSlaveStore(),
		Timer:  NewTimerStore(),
		Line:   NewLineStore(),
	}
	log.Sugar.Info("cache 空间生成....")
	return store
}
