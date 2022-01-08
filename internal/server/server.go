package server

import (
	"context"
	"fmt"
	"giot/internal/conf"
	"giot/internal/log"
	"github.com/xormplus/xorm"
	"os"

	"go.uber.org/zap"
	"net/http"
	"time"
)

type server struct {
	http *http.Server
	db   *xorm.Engine
	tcp  *TcpServer
}

func NewServer() *server {
	return &server{}
}

func (s *server) init() error {
	log.Info("Initialize postgres...")
	//err := s.setupDb()
	//if err != nil {
	//	return err
	//}
	log.Info("Initialize server...")
	s.setupServer()
	setupTcp()
	return nil
}

func (s *server) Start(er chan error) {
	err := s.init()
	if err != nil {
		er <- err
		return
	}
	s.printInfo()
	//start gin server
	log.Infof("start server Listening on: %s", s.http.Addr)
	go func() {
		err := s.http.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			log.Error("listen and serv fail: %s", zap.Error(err))
			er <- err
		}
	}()
}

func (s *server) Stop() {
	s.shutdownServer(s.http)

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
