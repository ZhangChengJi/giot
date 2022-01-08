package processor

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

type syncTimer struct {
	c          gnet.Conn
	guid       string
	directives [][]byte
}

type Interface interface {
	execute()
}

func (t *syncTimer) execute() {
	for _, v := range t.directives {
		t.c.AsyncWrite(v)
	}
}

type DeviceScheduler struct {
	Interval time.Duration
}

func (s *DeviceScheduler) Next(prev time.Time) time.Time {
	return prev.Add(s.Interval)
}

func scheduleTimer() {
	tw := timingwheel.NewTimingWheel(time.Millisecond, 20)
	tw.Start()
	defer tw.Stop()
	tw.ScheduleFunc(&DeviceScheduler{time.Second}, func() {

	})

}
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
