package server

import (
	"context"
	"fmt"
	"giot/conf"
	"giot/internal/scheduler/device"
	"giot/internal/scheduler/transfer"
	"giot/pkg/etcd"
	"giot/pkg/gorm"
	"giot/pkg/log"
	"giot/pkg/mqtt"
	"giot/pkg/tdengine"
	"math/rand"
	"os"
	"strconv"

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
}

func NewServer() *server {
	return &server{}
}

func (s *server) init() error {

	log.Info("Initialize mysql...")
	db, err := gorm.New(conf.MysqlConfig)
	if err != nil {
		fmt.Println(err)
		return err
	}
	log.Info("Initialize etcd...")
	err = etcd.InitETCDClient(conf.ETCDConfig)
	if err != nil {
		fmt.Println(err)
		return err
	}

	log.Info("Initialize tdengine...")
	td, err := tdengine.New(conf.TdengineConfig)
	if err != nil {
		fmt.Println(err)
		return err
	}

	log.Info("Initialize mqtt...")
	conf.MqttConfig.ClientId = "scheduler" + strconv.Itoa(rand.New(rand.NewSource(time.Now().UnixNano())).Int())

	mq, err := mqtt.New(conf.MqttConfig)
	if err != nil {
		fmt.Println(err)
		return err
	}

	err = device.Setup(etcd.GenEtcdStorage(), db, mq)
	if err != nil {
		fmt.Println(err)
		return err
	}

	go transfer.Setup(mq, td, db)

	return nil
}

func (s *server) Start(er chan error) {
	err := s.init()
	if err != nil {
		er <- err
		return
	}
	s.printInfo()
}

func (s *server) Stop() {
	s.shutdownServer(nil)

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
