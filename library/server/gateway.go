package server

import (
	"context"
	"fmt"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"google.golang.org/grpc"
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
	MuxOptions     []runtime.ServeMuxOption
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
		MuxOptions: []runtime.ServeMuxOption{
			runtime.WithMarshalerOption(runtime.MIMEWildcard, &runtime.JSONPb{EmitDefaults: true}),
			runtime.WithProtoErrorHandler(runtime.DefaultHTTPProtoErrorHandler),
		},
	}
}

func newGatewayServer(c *gatewayConfig, conn *grpc.ClientConn, servers []ServiceServer) (*gatewayServer, error) {
	mux := runtime.NewServeMux(c.MuxOptions...)
	httpMux := http.NewServeMux()
	httpPattern := "/http"

	for _, handler := range c.ServerHandlers {
		handler(httpMux)
	}
	//httpMux.Handle("/grpc/", mux)

	server := &http.Server{
		Addr:    c.Addr.String(),
		Handler: httpMux,
	}
	if cfg := c.ServerConfig; cfg != nil {
		cfg.applyTo(server)
	}

	for _, server := range servers {
		err := server.RegisterWithHandler(context.Background(), mux, conn)
		if err != nil {
			return nil, fmt.Errorf("Fail to register handler, %v\n", err)
		}
		handler, err := server.RegisterWithHttpHandler(httpPattern)
		if err != nil {
			return nil, fmt.Errorf("Fail to register http handler, %v\n", err)
		}
		httpMux.Handle(httpPattern+"/", handler)
	}

	return &gatewayServer{
		server: server,
		config: c,
	}, nil
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

func (c *HTTPServerConfig) applyTo(s *http.Server) {
	s.ReadTimeout = c.ReadTimeout
	s.WriteTimeout = c.WriteTimeout
	s.MaxHeaderBytes = c.MaxHeaderBytes
	s.ConnState = c.ConnState
}
