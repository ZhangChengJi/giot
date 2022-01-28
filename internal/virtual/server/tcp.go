package server

import (
	"context"
	"fmt"
	"giot/conf"
	"giot/internal/model"
	"giot/internal/virtual/device"
	"giot/internal/virtual/tcp"
	"github.com/panjf2000/gnet"
	"log"

	"time"
)

func (s *server) setupTcp() {
	config := conf.GnetConfig
	log.Println("gent tcp event loop started")
	t := &tcp.TcpServer{
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
	log.Fatalf("gent tcp event loop start failed: %v", gnet.Serve(t, fmt.Sprintf("tcp://%v", config.Addr), gnet.WithMulticore(config.Multicore), gnet.WithTCPKeepAlive(5*time.Second), gnet.WithReusePort(config.Reuseport)))
}

func (s *server) shutdown() error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	return gnet.Stop(ctx, conf.GnetConfig.Addr)
}
