package store

import (
	"context"
	"giot/pkg/log"
	cmap "github.com/orcaman/concurrent-map"
	"github.com/shiningrush/droplet/data"
)

type GuidStoreIn interface {
	Create(ctx context.Context, key string, obj string)
	Get(ctx context.Context, key string) (string, error)
	Update(ctx context.Context, key string, obj string) (string, error)
	Delete(ctx context.Context, key string)
	//DeleteTask(ctx context.Context, RemoteAddr string)
}

type GuidStore struct {
	cache cmap.ConcurrentMap
}

func NewGuidStore() *GuidStore {
	return &GuidStore{cache: cmap.New()}
}
func (t *GuidStore) Create(ctx context.Context, key string, obj string) {
	t.cache.Set(key, obj)
}
func (t *GuidStore) Get(ctx context.Context, key string) (string, error) {
	//key = strings.Trim(key, "\r")
	if tmp, ok := t.cache.Get(key); ok {
		return tmp.(string), nil
	} else {
		log.Warnf("data not found by key: %s", key)
		return "", data.ErrNotFound
	}
}
func (t *GuidStore) Update(ctx context.Context, key string, obj string) (string, error) {
	t.cache.Set(key, obj)
	return obj, nil
}
func (t *GuidStore) Delete(ctx context.Context, key string) {
	t.cache.Remove(key)
}
