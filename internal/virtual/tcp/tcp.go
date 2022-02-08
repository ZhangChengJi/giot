package tcp

import (
	"fmt"
	"giot/internal/model"
	modbus2 "giot/pkg/modbus"
	"github.com/panjf2000/gnet"
	"log"
	"sync"
	"time"
)

type TcpServer struct {
	*gnet.EventServer
	connectedSockets sync.Map
	codec            gnet.ICodec
	RegisterChan     chan *model.RegisterData
	DataChan         chan *model.RemoteData
	ListenMsgChan    chan *model.ListenMsg
}

func (ps *TcpServer) OnInitComplete(srv gnet.Server) (action gnet.Action) {
	log.Printf("Push server is listening on %s (multi-cores: %t, loops: %d)🚀...\n", srv.Addr.String(), srv.Multicore, srv.NumEventLoop)
	return
}

func (ps *TcpServer) OnOpened(c gnet.Conn) (out []byte, action gnet.Action) {
	fmt.Printf("新打开open:%X\n", c.RemoteAddr())
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

func (ps *TcpServer) React(frame []byte, c gnet.Conn) (out []byte, action gnet.Action) {

	fmt.Printf("上数data:%X\n", frame)

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
	out = []byte{0x01, 0x71, 0x00, 0x01, 0x00, 0x01, 0x00, 0x01}
	return
}

func (ps *TcpServer) Aaa() {
	fmt.Println("当前时间：", time.Now())
	myTicker := time.NewTicker(time.Second * 5) //
	go func() {
		for {
			<-myTicker.C
			r, _ := modbus2.NewClient(&modbus2.RtuHandler{}).WriteSingleRegister(1, 1, 1, modbus2.Success)
			fmt.Printf("%X", r)
			ps.connectedSockets.Range(func(key, value interface{}) bool {
				c := value.(gnet.Conn)
				c.AsyncWrite(r)
				return true
			})
		}
	}()
}
