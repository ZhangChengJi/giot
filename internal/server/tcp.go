package server

import (
	"fmt"
	"giot/internal/conf"
	"giot/internal/manager/modbus"
	"giot/internal/processor"

	"github.com/panjf2000/gnet"
	"log"
	"sync"
	"time"
)

type TcpServer struct {
	*gnet.EventServer
	connectedSockets sync.Map
	codec            gnet.ICodec
	data             processor.RemoteData
	re               processor.RegisterData
}

func (ps *TcpServer) OnInitComplete(srv gnet.Server) (action gnet.Action) {
	log.Printf("Push server is listening on %s (multi-cores: %t, loops: %d)🚀...\n", srv.Addr.String(), srv.Multicore, srv.NumEventLoop)
	return
}

func (ps *TcpServer) OnOpened(c gnet.Conn) (out []byte, action gnet.Action) {
	fmt.Printf("data:%X\n", c.Read())
	log.Printf("Socket with addr: %s has been opened...\n", c.RemoteAddr().String())
	ps.connectedSockets.Store(c.RemoteAddr().String(), c)
	return
}

func (ps *TcpServer) OnClosed(c gnet.Conn, err error) (action gnet.Action) {

	log.Printf("Socket with addr: %s is closing...\n", c.RemoteAddr().String())
	ps.connectedSockets.Delete(c.RemoteAddr().String())
	return
}

//func (ps *pushServer) Tick() (delay time.Duration, action gnet.Action) {
//	log.Println("It's time to push data to clients!!!")
//	ps.connectedSockets.Range(func(key, value interface{}) bool {
//		addr := key.(string)
//		c := value.(gnet.Conn)
//		c.AsyncWrite([]byte(fmt.Sprintf("heart beating to %s\n", addr)))
//		return true
//	})
//	delay = ps.tick
//	return
//}

func (ps *TcpServer) React(frame []byte, c gnet.Conn) (out []byte, action gnet.Action) {

	fmt.Printf("data:%X\n", frame)

	length := len(frame)
	if length > 0 && length > 7 {
		if length == 24 { //注册
			ps.re.C = c
			ps.re.D = frame
			processor.RegisterChan <- ps.re
		} else { //上数
			ps.data.RemoteIp = []byte(c.RemoteAddr().String())
			ps.data.Frame = frame
			processor.DataChan <- ps.data
		}
	}

	//data := append([]byte{}, frame...)
	//_ = ps.workerPool.Submit(func() {
	//	c.AsyncWrite(data)
	//})
	//	return
	//out = []byte{0x01, 0x71, 0x00, 0x01, 0x00, 0x01, 0x00, 0x01}
	return
}

func (ps *TcpServer) aaa() {
	fmt.Println("当前时间：", time.Now())
	myTicker := time.NewTicker(time.Second * 3) //
	go func() {
		for {
			<-myTicker.C
			r, _ := modbus.NewClient(&modbus.RtuHandler{}).ReadHoldingRegisters(0, 1)
			ps.connectedSockets.Range(func(key, value interface{}) bool {
				c := value.(gnet.Conn)
				c.AsyncWrite(r)
				return true
			})
		}
	}()
}

func setupTcp() {
	processor.Setup()
	config := conf.GnetConfig
	log.Println("gent tcp event loop started")
	t := &TcpServer{}
	t.aaa()
	log.Fatalf("gent tcp event loop start failed: %v", gnet.Serve(t, fmt.Sprintf("tcp://%v", config.Addr), gnet.WithMulticore(config.Multicore), gnet.WithReusePort(config.Reuseport)))
}
