package server

import (
	"context"
	"log"
	"net"
	"net/http"
	"time"
)

type gatewayServer struct {
	server *http.Server
	config *gatewayConfig
}

type gatewayConfig struct {
	Addr           Listen
	ServerConfig   *HTTPServerConfig
	ServerHandlers []HTTPServerHandler
}

type HTTPServerConfig struct {
	ReadTimeout    time.Duration
	WriteTimeout   time.Duration
	MaxHeaderBytes int
	ConnState      func(net.Conn, http.ConnState)
}

func createDefaultGatewayConfig() *gatewayConfig {
	return &gatewayConfig{
		Addr: Listen{
			Host: "0.0.0.0",
			Port: 10080,
		},
		ServerConfig:   nil,
		ServerHandlers: nil,
	}
}

func (s *gatewayServer) Serve() error {
	log.Printf("Http starting: %v", s.config.Addr.String())
	if err := s.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Printf("Error serve http server, %v\n", err)
		return err
	}
	return nil
}

func (s *gatewayServer) Shutdown(ctx context.Context) {
	log.Println("Http server is shutdown")
	err := s.server.Shutdown(ctx)
	if err != nil {
		log.Printf("Fail shut down gateway server, %v\n", err)
	}
}
