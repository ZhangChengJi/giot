package store

import (
	"context"
	"giot/internal/model"
	"giot/pkg/log"
	cmap "github.com/orcaman/concurrent-map"
	"github.com/shiningrush/droplet/data"
)

type SlaveStoreIn interface {
	Create(ctx context.Context, key string, obj []*model.Slave)
	Get(ctx context.Context, key string) ([]*model.Slave, error)
	GetSlave(ctx context.Context, key string, slaveId byte) (*model.Slave, error)
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
		log.Sugar.Warnf("data not found by key: %s", key)
		return nil, data.ErrNotFound
	}
}
func (t *SlaveStore) GetSlave(ctx context.Context, key string, slaveId byte) (*model.Slave, error) {
	s, err := t.Get(ctx, key)
	if err != nil {
		return nil, err
	}

	for _, a := range s {
		if slaveId == a.SlaveId {
			return a, nil
		}
	}
	log.Sugar.Warnf("attributeId not found by key: %s", key)
	return nil, data.ErrNotFound
}
func (t *SlaveStore) Update(ctx context.Context, key string, obj []*model.Slave) ([]*model.Slave, error) {
	t.cache.Set(key, obj)
	return obj, nil
}
func (t *SlaveStore) Delete(ctx context.Context, key string) {
	t.cache.Remove(key)
}
