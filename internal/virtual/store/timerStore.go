package store

import (
	"context"
	"fmt"
	"giot/internal/virtual/wheelTimer"
	"giot/pkg/log"
	cmap "github.com/orcaman/concurrent-map"
	"github.com/panjf2000/gnet/v2"
	"github.com/shiningrush/droplet/data"
)

type DeviceTimerIn interface {
	Create(ctx context.Context, key string, obj *wheelTimer.SyncTimer)
	Get(ctx context.Context, key string) (*wheelTimer.SyncTimer, error)
	GetConn(ctx context.Context, key string) (gnet.Conn, error)
	Update(ctx context.Context, key string, obj *wheelTimer.SyncTimer) (*wheelTimer.SyncTimer, error)
	Delete(ctx context.Context, key string)
	//DeleteTask(ctx context.Context, RemoteAddr string)
}

type TimerStore struct {
	cache cmap.ConcurrentMap
	//guid  cmap.ConcurrentMap
}

func NewTimerStore() *TimerStore {
	return &TimerStore{cache: cmap.New()}
}

func (t *TimerStore) Create(ctx context.Context, key string, obj *wheelTimer.SyncTimer) {
	t.cache.Set(key, obj)
	//	t.guid.Set(RemoteAddr, key)
}
func (t *TimerStore) Get(ctx context.Context, key string) (*wheelTimer.SyncTimer, error) {
	fmt.Println("key:", t.cache.Keys())
	if tmp, ok := t.cache.Get(key); ok {
		return tmp.(*wheelTimer.SyncTimer), nil
	} else {
		log.Sugar.Warnf("data not found by key: %s", key)
		return nil, data.ErrNotFound
	}
}
func (t *TimerStore) Update(ctx context.Context, key string, obj *wheelTimer.SyncTimer) (*wheelTimer.SyncTimer, error) {
	t.cache.Set(key, obj)
	return obj, nil
}
func (t *TimerStore) Delete(ctx context.Context, key string) {
	t.cache.Remove(key)

}
func (t *TimerStore) GetConn(ctx context.Context, key string) (gnet.Conn, error) {
	timers, err := t.Get(ctx, key)
	if err != nil {
		return nil, err
	}
	return timers.Conn, nil

}

//func (t *TimerStore) DeleteTask(ctx context.Context, RemoteAddr string) {
//	if id, ok := t.guid.Get(RemoteAddr); ok {
//		tw, err := t.Get(context.TODO(), id.(string))
//		if err != nil {
//			logs.Warnf("RemoteAddr:%s task not empty", RemoteAddr)
//			return
//		}
//		tw.T.Stop()
//		tw.Conn.Close()
//		fmt.Printf("删除此任务:%s", tw)
//		t.cache.Remove(id.(string))
//		t.guid.Remove(RemoteAddr)
//
//	}
//
//}
