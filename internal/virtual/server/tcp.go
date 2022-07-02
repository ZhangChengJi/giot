package server

import (
	"context"
	"fmt"
	"giot/conf"
	"giot/internal/model"
	"giot/internal/virtual/device"
	"giot/internal/virtual/tcp"
	"giot/pkg/log"
	"github.com/panjf2000/gnet/v2"

	"time"
)

func (s *server) setupTcp() {
	config := conf.GnetConfig
	log.Sugar.Infof("gent tcp event loop started...🚀 🚀 🚀")
	t := &tcp.TcpServer{
		DataChan:      make(chan *model.RemoteData),
		RegisterChan:  make(chan *model.RegisterData),
		ListenMsgChan: make(chan *model.ListenMsg),
	}
	pro := NewProcessor()
	pro.Timer.Start()
	defer pro.Timer.Stop()
	device.Init()
	go pro.Swift(t.RegisterChan)
	go pro.Handle(t.DataChan)
	go pro.ListenCommand(t.ListenMsgChan)
	log.Sugar.Fatalf("gent tcp event loop start failed: %v", gnet.Run(t, fmt.Sprintf("tcp://%v", config.Addr), gnet.WithMulticore(config.Multicore), gnet.WithTCPKeepAlive(5*time.Second), gnet.WithReusePort(config.Reuseport)))
}

func (s *server) shutdown() error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	return gnet.Stop(ctx, conf.GnetConfig.Addr)
}
