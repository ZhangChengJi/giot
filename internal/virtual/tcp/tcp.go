package tcp

import (
	"giot/internal/model"
	"giot/internal/virtual/protocol"
	"giot/pkg/log"
	"github.com/panjf2000/gnet"
	"strings"
	"sync"
)

var (
	ListenMsgChan chan *model.ListenMsg
)

type TcpServer struct {
	*gnet.EventServer
	connectedSockets sync.Map
	RegisterChan     chan *model.RegisterData
	DataChan         chan *model.RemoteData
}

func (ps *TcpServer) OnOpened(c gnet.Conn) (out []byte, action gnet.Action) {
	c.SetContext(new(protocol.ModbusCodec))
	log.Sugar.Infof("Socket with addr: %s has been opened...\n", c.RemoteAddr().String())
	ps.connectedSockets.Store(c.RemoteAddr().String(), c)
	return
}

func (ps *TcpServer) OnClosed(c gnet.Conn, err error) (action gnet.Action) {
	log.Sugar.Infof("Socket with addr: %s is closing...\n", c.RemoteAddr().String())
	ps.connectedSockets.Delete(c.RemoteAddr().String())
	ha := &model.ListenMsg{
		ListenType: 1,
		RemoteAddr: c.RemoteAddr().String(),
		Command:    1,
	}
	ListenMsgChan <- ha
	return
}
func (ps *TcpServer) React(frame []byte, c gnet.Conn) (out []byte, action gnet.Action) {
	data := make([]byte, len(frame))
	copy(data, frame)
	length := len(data)
	if length > 0 && length >= 7 && length <= 44 {
		if length == 24 || length == 44 { //注册
			re := &model.RegisterData{Conn: c, D: strings.Trim(string(data), "\r")}
			ps.RegisterChan <- re
		} else { //上数
			da := &model.RemoteData{
				Frame:      data,
				RemoteAddr: c.RemoteAddr().String(),
				Conn:       c,
			}
			ps.DataChan <- da
		}
	}
	return
}
