package lineTimer

import (
	"context"
	"fmt"
	"giot/internal/virtual/device"
	"giot/internal/virtual/store"
	"giot/pkg/log"
	"giot/utils/consts"
	"github.com/RussellLuo/timingwheel"
	cmap "github.com/orcaman/concurrent-map"
	"github.com/shiningrush/droplet/data"
	"time"
)

type LineTimer struct {
	RemoteAddr string
	Guid       string
	Time       time.Duration
	SlaveStore store.SlaveStoreIn
	T          *timingwheel.Timer
}

func NewTimer() *timingwheel.TimingWheel {
	return scheduleTimer()
}

type DeviceScheduler struct {
	Interval time.Duration
}

func (s *DeviceScheduler) Next(prev time.Time) time.Time {
	return prev.Add(s.Interval)
}

func scheduleTimer() *timingwheel.TimingWheel {
	tw := timingwheel.NewTimingWheel(time.Second, 20)
	return tw
}

func (t *LineTimer) Execute() {
	list, err := t.SlaveStore.Get(context.TODO(), t.RemoteAddr)
	if err != nil {
		return
	}
	for _, slave := range list {
		if time.Now().Sub(slave.DataTime) > t.Time { //下线
			fmt.Printf("时间:%v----->slave:%v下⬇️线\n", time.Now().Format("2006-01-02 15:04:05"), slave.SlaveId)
			if slave.LineStatus == "" || slave.LineStatus == consts.ONLINE {
				fmt.Println("进入。。。。。")
				slave.LineStatus = consts.OFFLINE
				device.OnlineChan <- &device.DeviceMsg{Ts: time.Now(), Status: consts.OFFLINE, DeviceId: t.Guid, SlaveId: int(slave.SlaveId)}
			}
		} else {
			fmt.Printf("时间:%v----->slave:%v上⬆️线\n", time.Now().Format("2006-01-02 15:04:05"), slave.SlaveId)
			if slave.LineStatus == "" || slave.LineStatus == consts.OFFLINE { //如果之前有下线的现在变为上线
				slave.LineStatus = consts.ONLINE
				device.OnlineChan <- &device.DeviceMsg{Ts: time.Now(), Status: consts.ONLINE, DeviceId: t.Guid, SlaveId: int(slave.SlaveId)}
			}
		}
	}

}

type LineStoreIn interface {
	Create(ctx context.Context, key string, obj *LineTimer)
	Get(ctx context.Context, key string) (*LineTimer, error)
	Update(ctx context.Context, key string, obj *LineTimer) (*LineTimer, error)
	Delete(ctx context.Context, key string)
	//DeleteTask(ctx context.Context, RemoteAddr string)
}

type LineStore struct {
	cache cmap.ConcurrentMap
	//guid  cmap.ConcurrentMap
}

func NewLineStore() *LineStore {
	return &LineStore{cache: cmap.New()}
}

func (t *LineStore) Create(ctx context.Context, key string, obj *LineTimer) {
	t.cache.Set(key, obj)
	//	t.guid.Set(RemoteAddr, key)
}
func (t *LineStore) Get(ctx context.Context, key string) (*LineTimer, error) {
	fmt.Println("key:", t.cache.Keys())
	if tmp, ok := t.cache.Get(key); ok {
		return tmp.(*LineTimer), nil
	} else {
		log.Warnf("data not found by key: %s", key)
		return nil, data.ErrNotFound
	}
}
func (t *LineStore) Update(ctx context.Context, key string, obj *LineTimer) (*LineTimer, error) {
	t.cache.Set(key, obj)
	return obj, nil
}
func (t *LineStore) Delete(ctx context.Context, key string) {
	t.cache.Remove(key)

}

//func (t *TimerStore) DeleteTask(ctx context.Context, RemoteAddr string) {
//	if id, ok := t.guid.Get(RemoteAddr); ok {
//		tw, err := t.Get(context.TODO(), id.(string))
//		if err != nil {
//			log.Warnf("RemoteAddr:%s task not empty", RemoteAddr)
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
