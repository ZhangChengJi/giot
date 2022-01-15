package store

import (
	"context"
	"giot/internal/core/model"
	"giot/internal/log"
	cmap "github.com/orcaman/concurrent-map"
	"github.com/shiningrush/droplet/data"
)

type SlaveStoreIn interface {
	Create(ctx context.Context, key string, obj []*model.Slave)
	Get(ctx context.Context, key string) ([]*model.Slave, error)
	GetAttributeId(ctx context.Context, key string, slaveId byte) (string, string, error)
	Update(ctx context.Context, key string, obj []*model.Slave) ([]*model.Slave, error)
	Delete(ctx context.Context, key string)
	//DeleteTask(ctx context.Context, RemoteAddr string)
}

type SlaveStore struct {
	cache cmap.ConcurrentMap
}

func NewSlaveStore() *SlaveStore {
	return &SlaveStore{cache: cmap.New()}
}
func (t *SlaveStore) Create(ctx context.Context, key string, obj []*model.Slave) {
	t.cache.Set(key, obj)
}
func (t *SlaveStore) Get(ctx context.Context, key string) ([]*model.Slave, error) {
	if tmp, ok := t.cache.Get(key); ok {
		return tmp.([]*model.Slave), nil
	} else {
		log.Warnf("data not found by key: %s", key)
		return nil, data.ErrNotFound
	}
}
func (t *SlaveStore) GetAttributeId(ctx context.Context, key string, slaveId byte) (string, string, error) {
	s, err := t.Get(ctx, key)
	if err != nil {
		return "", "", err
	}

	for _, a := range s {
		if slaveId == a.SlaveId {
			return a.AttributeId, a.DeviceId, nil
		}
	}
	log.Warnf("attributeId not found by key: %s", key)
	return "", "", data.ErrNotFound
}
func (t *SlaveStore) Update(ctx context.Context, key string, obj []*model.Slave) ([]*model.Slave, error) {
	t.cache.Set(key, obj)
	return obj, nil
}
func (t *SlaveStore) Delete(ctx context.Context, key string) {
	t.cache.Remove(key)
}
