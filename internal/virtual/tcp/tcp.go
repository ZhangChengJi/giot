package tcp

import (
	"fmt"
	"giot/internal/model"
	"giot/internal/virtual/protocol"
	"giot/pkg/log"
	"github.com/panjf2000/gnet/v2"
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
	log.Sugar.Infof("Socket with addr: %s has been opened...\n", c.RemoteAddr().String())
	ps.connectedSockets.Store(c.RemoteAddr().String(), c)
	return
}

func (ps *TcpServer) OnClose(c gnet.Conn, err error) (action gnet.Action) {
	log.Sugar.Infof("Socket with addr: %s is closing...\n", c.RemoteAddr().String())
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
	length := len(data)
	if length > 0 && length >= 7 && length < 44 {
		if length == 24 { //注册
			re := &model.RegisterData{Conn: c, D: strings.Trim(string(data), "\r")}
			ps.RegisterChan <- re
		} else { //上数
			da := &model.RemoteData{
				Frame:      data,
				RemoteAddr: c.RemoteAddr().String(),
				Conn:       c,
			}
			fmt.Println(da)
			ps.DataChan <- da
		}
	}
	return gnet.None
}
