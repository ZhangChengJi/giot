package offline

import (
	"context"
	"fmt"
	"giot/internal/model"
	"giot/internal/virtual/device"
	"giot/internal/virtual/mqtt"
	"giot/internal/virtual/store"
	"giot/pkg/log"
	"giot/utils/consts"
	"giot/utils/json"
	"github.com/RussellLuo/timingwheel"
	cmap "github.com/orcaman/concurrent-map"
	"github.com/shiningrush/droplet/data"
	"time"
)

func (t *LineTimer) deleteTask(remoteAddr string) {
	ha := &model.ListenMsg{
		ListenType: 1,
		RemoteAddr: remoteAddr,
		Command:    1,
	}
	t.DeleteMsg <- ha
}
func (t *LineTimer) line(data *device.DeviceMsg) {
	topic := append([]byte("device/line/"), data.DeviceId...)
	payload, _ := json.Marshal(data)
	t.MqttBroker.Publish(string(topic), payload)
}

type LineTimer struct {
	RemoteAddr string
	Guid       string
	Time       time.Duration
	Status     int
	initOp     int
	noDataNum  int
	GuidStore  store.GuidStoreIn
	SlaveStore store.SlaveStoreIn
	MqttBroker mqtt.Broker
	T          *timingwheel.Timer
	DeleteMsg  func()
}

func NewTimer() *timingwheel.TimingWheel {
	return scheduleTimer()
}

type OfflineScheduler struct {
	Interval time.Duration
}

func (s *OfflineScheduler) Next(prev time.Time) time.Time {
	return prev.Add(s.Interval)
}

func scheduleTimer() *timingwheel.TimingWheel {
	tw := timingwheel.NewTimingWheel(time.Second, 3600)
	return tw
}

func (t *LineTimer) Execute() {
	if address, err := t.GuidStore.Get(context.TODO(), t.Guid); err == nil && address != "" {
		var i int
		if t.Status > 0 { //之前是在线
			slaveList, err := t.SlaveStore.Get(context.TODO(), address)
			if err != nil {
				return
			}

			for _, slave := range slaveList {
				if time.Now().Sub(slave.DataTime) > t.Time { //下线
					fmt.Printf("时间:%v----->slave:%v下⬇️线\n", time.Now().Format("2006-01-02 15:04:05"), slave.SlaveId)
					if slave.LineStatus == "" || slave.LineStatus == consts.ONLINE {
						fmt.Println("进入。。。。。")
						slave.LineStatus = consts.OFFLINE
						i++
						t.line(&device.DeviceMsg{Ts: time.Now(), Status: consts.OFFLINE, DeviceId: t.Guid, SlaveId: int(slave.SlaveId)})
					}
				} else { //上线
					t.noDataNum = 0
					fmt.Printf("时间:%v----->slave:%v上⬆️线\n", time.Now().Format("2006-01-02 15:04:05"), slave.SlaveId)
					if slave.LineStatus == "" || slave.LineStatus == consts.OFFLINE { //如果之前有下线的现在变为上线
						fmt.Println("进入。。。。。")
						slave.LineStatus = consts.ONLINE
						t.line(&device.DeviceMsg{Ts: time.Now(), Status: consts.ONLINE, DeviceId: t.Guid, SlaveId: int(slave.SlaveId)})
					}
				}
			}
			if len(slaveList) >= i {
				t.noDataNum++
			}
			if t.noDataNum >= 2 {
				t.deleteTask(address)
				t.Status = 0 //改为离线
				t.line(&device.DeviceMsg{Ts: time.Now(), Status: consts.OFFLINE, DeviceId: t.Guid})
			}

		} else { //之前是离线
			t.Status = 1 //改为上线
			slaveList, err := t.SlaveStore.Get(context.TODO(), t.Guid)
			if err != nil {
				return
			}
			for _, slave := range slaveList {
				if time.Now().Sub(slave.DataTime) > t.Time { //下线
					fmt.Printf("时间:%v----->slave:%v下⬇️线\n", time.Now().Format("2006-01-02 15:04:05"), slave.SlaveId)
					if slave.LineStatus == "" || slave.LineStatus == consts.ONLINE {
						fmt.Println("进入。。。。。")
						slave.LineStatus = consts.OFFLINE
						i++
						t.line(&device.DeviceMsg{Ts: time.Now(), Status: consts.OFFLINE, DeviceId: t.Guid, SlaveId: int(slave.SlaveId)})
					}
				} else {
					t.noDataNum = 0
					fmt.Printf("时间:%v----->slave:%v上⬆️线\n", time.Now().Format("2006-01-02 15:04:05"), slave.SlaveId)
					if slave.LineStatus == "" || slave.LineStatus == consts.OFFLINE { //如果之前有下线的现在变为上线
						fmt.Println("进入。。。。。")
						slave.LineStatus = consts.ONLINE
						t.line(&device.DeviceMsg{Ts: time.Now(), Status: consts.ONLINE, DeviceId: t.Guid, SlaveId: int(slave.SlaveId)})
					}
				}
			}
			t.line(&device.DeviceMsg{Ts: time.Now(), Status: consts.ONLINE, DeviceId: t.Guid})
			if len(slaveList) >= i {
				t.noDataNum++
			}
			if t.noDataNum >= 2 {
				t.deleteTask(address)
				t.Status = 0 //改为离线
				t.line(&device.DeviceMsg{Ts: time.Now(), Status: consts.OFFLINE, DeviceId: t.Guid})
			}
		}
	} else if address == "" { //没有连接情况下
		if t.initOp == 0 { //之前是在线
			t.Status = 0 //改为离线
			t.initOp = 1
			t.line(&device.DeviceMsg{Ts: time.Now(), Status: consts.OFFLINE, DeviceId: t.Guid})

			//} else { //之前是离线

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
		log.Sugar.Warnf("data not found by key: %s", key)
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
func metaDataCompile(val string) (*model.Device, error) {
	devic := &model.Device{}
	err := json.Unmarshal([]byte(val), devic)
	if err != nil {
		log.Sugar.Errorf("json unmarshal failed: %s", err)
		return nil, err
	}
	return devic, nil
}
