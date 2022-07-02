package store

import (
	"context"
	"giot/internal/model"
	"giot/pkg/log"
	cmap "github.com/orcaman/concurrent-map"
	"github.com/shiningrush/droplet/data"
)

type DeviceStoreIn interface {
	Create(ctx context.Context, key string, obj *model.Device)
	Get(ctx context.Context, key string) (*model.Device, error)
	//GetSlave(ctx context.Context, key string, slaveId byte) (*model.Slave, error)
	Update(ctx context.Context, key string, obj *model.Device) (*model.Device, error)
	Delete(ctx context.Context, key string)
	//DeleteTask(ctx context.Context, RemoteAddr string)
}

type DeviceStore struct {
	cache cmap.ConcurrentMap
}

func NewDeviceStore() *DeviceStore {
	return &DeviceStore{cache: cmap.New()}
}
func (t *DeviceStore) Create(ctx context.Context, key string, obj *model.Device) {
	t.cache.Set(key, obj)
}
func (t *DeviceStore) Get(ctx context.Context, key string) (*model.Device, error) {
	if tmp, ok := t.cache.Get(key); ok {
		return tmp.(*model.Device), nil
	} else {
		log.Sugar.Warnf("data not found by key: %s", key)
		return nil, data.ErrNotFound
	}
}

//func (t *DeviceStore) GetSlave(ctx context.Context, key string, slaveId byte) (*model.Slave, error) {
//	s, err := t.Get(ctx, key)
//	if err != nil {
//		return nil, err
//	}
//
//	for _, a := range s {
//		if slaveId == a.SlaveId {
//			return a, nil
//		}
//	}
//	logs.Warnf("attributeId not found by key: %s", key)
//	return nil, data.ErrNotFound
//}
func (t *DeviceStore) Update(ctx context.Context, key string, obj *model.Device) (*model.Device, error) {
	t.cache.Set(key, obj)
	return obj, nil
}
func (t *DeviceStore) Delete(ctx context.Context, key string) {
	t.cache.Remove(key)
}
