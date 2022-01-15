package server

import (
	"fmt"
	"giot/internal/conf"
	"giot/internal/core/model"
	"giot/internal/device"
	"giot/internal/manager/modbus"
	"github.com/panjf2000/gnet"
	"log"
	"sync"
	"time"
)

type TcpServer struct {
	*gnet.EventServer
	connectedSockets sync.Map
	codec            gnet.ICodec
	//data             processor.RemoteData
	//re               processor.RegisterData
	RegisterChan  chan *model.RegisterData
	DataChan      chan *model.RemoteData
	ListenMsgChan chan *model.ListenMsg
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
	ha := &model.ListenMsg{
		ListenType: 1,
		RemoteAddr: c.RemoteAddr().String(),
		Command:    1,
	}
	ps.ListenMsgChan <- ha
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
	if length > 0 && length >= 7 {
		if length == 24 { //注册
			re := &model.RegisterData{C: c, D: frame}
			ps.RegisterChan <- re
		} else { //上数
			da := &model.RemoteData{
				Frame:      frame,
				RemoteAddr: c.RemoteAddr().String(),
			}
			ps.DataChan <- da
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
	myTicker := time.NewTicker(time.Second * 30) //
	go func() {
		for {
			<-myTicker.C
			r, _ := modbus.NewClient(&modbus.RtuHandler{}).ReadHoldingRegisters(1, 0, 1)
			fmt.Printf("%X", r)
			ps.connectedSockets.Range(func(key, value interface{}) bool {
				c := value.(gnet.Conn)
				c.AsyncWrite(r)
				return true
			})
		}
	}()
}

func (s *server) setupTcp() {
	config := conf.GnetConfig
	log.Println("gent tcp event loop started")
	t := &TcpServer{
		DataChan:      make(chan *model.RemoteData, 1024),
		RegisterChan:  make(chan *model.RegisterData, 1024),
		ListenMsgChan: make(chan *model.ListenMsg),
	}
	pro := NewProcessor()
	pro.Tw.Start()
	defer pro.Tw.Stop()
	go pro.Swift(t.DataChan, t.RegisterChan)
	go pro.ListenCommand(t.ListenMsgChan)
	device.Init()
	//	t.aaa()
	log.Fatalf("gent tcp event loop start failed: %v", gnet.Serve(t, fmt.Sprintf("tcp://%v", config.Addr), gnet.WithMulticore(config.Multicore), gnet.WithTCPKeepAlive(5*time.Second), gnet.WithReusePort(config.Reuseport)))
}
