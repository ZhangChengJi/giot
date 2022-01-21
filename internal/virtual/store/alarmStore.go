package store

import (
	"context"
	"giot/internal/virtual/engine"
	"giot/pkg/log"
	cmap "github.com/orcaman/concurrent-map"
	"github.com/shiningrush/droplet/data"
)

type AlarmStoreIn interface {
	Create(ctx context.Context, key string, obj engine.Interface)
	Get(ctx context.Context, key string) (engine.Interface, error)
	Update(ctx context.Context, key string, obj engine.Interface) (engine.Interface, error)
	Delete(ctx context.Context, key string)
	//DeleteTask(ctx context.Context, RemoteAddr string)
}

type AlarmStore struct {
	cache cmap.ConcurrentMap
}

func NewAlarmStore() *AlarmStore {
	return &AlarmStore{cache: cmap.New()}
}
func (t *AlarmStore) Create(ctx context.Context, key string, obj engine.Interface) {
	t.cache.Set(key, obj)
}
func (t *AlarmStore) Get(ctx context.Context, key string) (engine.Interface, error) {
	if tmp, ok := t.cache.Get(key); ok {
		return tmp.(engine.Interface), nil
	} else {
		log.Warnf("data not found by key: %s", key)
		return nil, data.ErrNotFound
	}
}
func (t *AlarmStore) Update(ctx context.Context, key string, obj engine.Interface) (engine.Interface, error) {
	t.cache.Set(key, obj)
	return obj, nil
}
func (t *AlarmStore) Delete(ctx context.Context, key string) {
	t.cache.Remove(key)
}
