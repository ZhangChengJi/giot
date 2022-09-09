package wheelTimer

import (
	"fmt"
	"github.com/RussellLuo/timingwheel"
	"github.com/panjf2000/gnet"
	"github.com/panjf2000/gnet/pkg/pool/goroutine"
	"sync"
	"time"
)

var (
	rtr        sync.Map
	rtp        sync.Map
	workerPool *goroutine.Pool
)

type SyncTimer struct {
	Conn       gnet.Conn
	RemoteAddr string
	Guid       string
	Time       time.Duration
	Directives [][]byte //下发指令
	T          *timingwheel.Timer
}

type Interface interface {
	Execute()
}

func (t *SyncTimer) Execute() {
	for _, v := range t.Directives {
		fmt.Printf("时间:%v——--->任务下发:%X\n", time.Now().Format("2006-01-02 15:04:05"), v)

		t.Conn.AsyncWrite(v)
		time.Sleep(300 * time.Millisecond)
	}
}

// DeviceScheduler **********************定时器***************************//

func NewTimer() *timingwheel.TimingWheel {
	return scheduleTimer()
}

type DeviceScheduler struct {
	Interval time.Duration
	current  int
	Rew      string
}

func (s *DeviceScheduler) Next(prev time.Time) time.Time {
	fmt.Printf("定时任务:%v:%v--->%v\n", s.Rew, s.current, s.Interval)
	s.current += 1
	return prev.Add(s.Interval)
}

func scheduleTimer() *timingwheel.TimingWheel {
	tw := timingwheel.NewTimingWheel(time.Second, 3600)
	return tw
}

//*********************END****************************//

//func Second60Timer() {
//	workerPool.Submit(func() {
//		t := time.NewTicker(time.Second * 60)
//		for {
//			//
//			rtp.Range(func(key, value interface{}) bool {
//				addr := key.(string)
//				c := value.(gnet.Conn)
//				c.AsyncWrite([]byte(fmt.Sprintf("heart beating to %s\n", addr)))
//				return true
//			})
//			<-t.C
//
//		}
//	})
//}
