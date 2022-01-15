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
		fmt.Println("任务下发")
		t.Conn.AsyncWrite(v)
	}
}

// DeviceScheduler **********************定时器***************************//

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

//*********************END****************************//

func Second60Timer() {
	workerPool.Submit(func() {
		t := time.NewTicker(time.Second * 60)
		for {
			//
			rtp.Range(func(key, value interface{}) bool {
				addr := key.(string)
				c := value.(gnet.Conn)
				c.AsyncWrite([]byte(fmt.Sprintf("heart beating to %s\n", addr)))
				return true
			})
			<-t.C

		}
	})
}
