package server

import (
	"context"
	"giot/internal/log"

	"fmt"
	"go.uber.org/zap"
	"net/http"
	"time"
)

const (
	systemAddress = 8080
)

type server struct {
	http *http.Server
}

type name interface {
}

func (s *server) init() error {

	log.Info("Initialize server...")
	s.setupServer()
	return nil
}

func (s *server) Start(er chan error) {
	err := s.init()
	if err != nil {
		er <- err
		return
	}
	//start gin server
	log.Info("start server Listening port: %s", zap.String("", s.http.Addr))
	go func() {
		err := s.http.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			log.Error("listen and serv fail: %s", zap.Error(err))
			er <- err
		}
	}()
}
func NewServer() *server {
	return &server{}
}

func (s *server) setupServer() {
	address := fmt.Sprintf(":%d", systemAddress)
	r := routers()
	s.http = &http.Server{
		Addr:           address,
		Handler:        r,
		ReadTimeout:    time.Duration(1000) * time.Millisecond,
		WriteTimeout:   time.Duration(5000) * time.Millisecond,
		MaxHeaderBytes: 1 << 20,
	}
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
