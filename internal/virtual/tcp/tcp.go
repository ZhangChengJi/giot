package tcp

import (
	"fmt"
	"giot/internal/model"
	"giot/internal/virtual/protocol"
	"github.com/panjf2000/gnet/v2"
	"log"
	"strings"
	"sync"
)

type TcpServer struct {
	gnet.BuiltinEventEngine

	connectedSockets sync.Map
	RegisterChan     chan *model.RegisterData
	DataChan         chan *model.RemoteData
	ListenMsgChan    chan *model.ListenMsg
}

func (ps *TcpServer) OnOpen(c gnet.Conn) (out []byte, action gnet.Action) {
	c.SetContext(new(protocol.ModbusCodec))
	log.Printf("Socket with addr: %s has been opened...\n", c.RemoteAddr().String())
	ps.connectedSockets.Store(c.RemoteAddr().String(), c)
	return
}

func (ps *TcpServer) OnClose(c gnet.Conn, err error) (action gnet.Action) {
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
func (ps *TcpServer) OnTraffic(c gnet.Conn) gnet.Action {

	data, _ := c.Next(-1)
	//c.AsyncWrite(data, nil)
	fmt.Println("data:", data)
	length := len(data)
	if length > 0 && length >= 7 && length < 30 {
		if length == 24 { //注册
			re := &model.RegisterData{Conn: c, D: strings.Trim(string(data), "\r")}
			ps.RegisterChan <- re
		} else { //上数
			da := &model.RemoteData{
				Frame:      data,
				RemoteAddr: c.RemoteAddr().String(),
			}
			ps.DataChan <- da
		}
	}
	return gnet.None
}

//func (ps *TcpServer) React(frame []byte, c gnet.Conn) (out []byte, action gnet.Action) {
//	fmt.Println("data:", frame)
//	length := len(frame)
//	if length > 0 && length >= 7 && length < 30 {
//		if length == 24 { //注册
//			re := &model.RegisterData{Conn: c, D: strings.Trim(string(frame), "\r")}
//			ps.RegisterChan <- re
//		} else { //上数
//			da := &model.RemoteData{
//				Frame:      frame,
//				RemoteAddr: c.RemoteAddr().String(),
//			}
//			ps.DataChan <- da
//		}
//	}
//	return
//}

//func (ps *TcpServer) Aaa() {
//	fmt.Println("当前时间：", time.Now())
//	myTicker := time.NewTicker(time.Second * 5) //
//	go func() {
//		for {
//			<-myTicker.C
//			r, _ := modbus2.NewClient(&modbus2.RtuHandler{}).WriteSingleRegister(1, 1, 1, modbus2.Success)
//			fmt.Printf("%X", r)
//			ps.connectedSockets.Range(func(key, value interface{}) bool {
//				c := value.(gnet.Conn)
//				c.AsyncWrite(r, nil)
//				return true
//			})
//		}
//	}()
//}
