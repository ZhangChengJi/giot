package server

import (
	"context"
	"fmt"
	"giot/conf"
	"giot/internal/virtual/tcp"
	"giot/pkg/log"
	"github.com/xormplus/xorm"
	"os"

	"go.uber.org/zap"
	"net/http"
	"time"
)

const (
	TCP = iota
	STORE
	NOTIFY
)

type server struct {
	db  *xorm.Engine
	tcp *tcp.TcpServer
}

func NewServer() *server {
	return &server{}
}

func (s *server) init() error {

	log.Info("Initialize etcd...")
	err := s.setupStore()
	if err != nil {
		fmt.Println(err)
		return err
	}
	log.Info("Initialize mqtt...")
	err = s.setupMqtt()
	if err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}

func (s *server) Start(er chan error) {
	err := s.init()
	if err != nil {
		er <- err
		return
	}
	go s.setupTcp()
	s.printInfo()
}

func (s *server) Stop() {
	s.shutdown()

}

func (s *server) shutdownServer(server *http.Server) {
	if server != nil {
		ctx, cancel := context.WithTimeout(context.TODO(), 5*time.Second)
		defer cancel()

		if err := server.Shutdown(ctx); err != nil {
			log.Error("Shutting down server error: %s", zap.Error(err))
		}
	}
}
func (s *server) printInfo() {
	fmt.Fprint(os.Stdout, "The giot is running successfully!\n\n")
	//utils.PrintVersion()
	fmt.Fprintf(os.Stdout, "%-8s: %s:%d\n", "Listen", conf.ServerHost, conf.ServerPort)
	if conf.SSLCert != "" && conf.SSLKey != "" {
		fmt.Fprintf(os.Stdout, "%-8s: %s:%d\n", "HTTPS Listen", conf.SSLHost, conf.SSLPort)
	}
	fmt.Fprintf(os.Stdout, "%-8s: %s\n", "Loglevel", conf.ErrorLogLevel)
	fmt.Fprintf(os.Stdout, "%-8s: %s\n", "ErrorLogFile", conf.ErrorLogPath)
	fmt.Fprintf(os.Stdout, "%-8s: %s\n\n", "AccessLogFile", conf.AccessLogPath)
}