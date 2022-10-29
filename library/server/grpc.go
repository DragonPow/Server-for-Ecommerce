package server

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"log"
)

type grpcServer struct {
	server *grpc.Server
	config *grpcConfig
}

type grpcConfig struct {
	Addr                Listen
	ServerOption        []grpc.ServerOption
	MaxConcurrentStream uint32
}

func createDefaultGrpcConfig() *grpcConfig {
	return &grpcConfig{
		Addr: Listen{
			Host: "0.0.0.0",
			Port: 10443,
		},
		ServerOption:        nil,
		MaxConcurrentStream: 1000,
	}
}

func (c grpcConfig) ServerOptions() []grpc.ServerOption {
	return append([]grpc.ServerOption{
		grpc.MaxConcurrentStreams(c.MaxConcurrentStream),
	}, c.ServerOption...)
}

func newGrpcServer(c *grpcConfig, servers []ServiceServer) *grpcServer {
	s := grpc.NewServer(c.ServerOptions()...)
	for _, server := range servers {
		server.RegisterWithServer(s)
	}
	return &grpcServer{
		server: s,
		config: c,
	}
}

func (s *grpcServer) Serve() error {
	l, err := s.config.Addr.CreateListener()
	if err != nil {
		return fmt.Errorf("Grpc fail create Listener, %v\n", err)
	}
	log.Printf("Starting grpc server, %v", l.Addr())
	err = s.server.Serve(l)
	if err != nil {
		log.Println(err)
		return fmt.Errorf("Error serve grpc server, %v\n", err)
	}
	return nil
}

func (s *grpcServer) Shutdown(ctx context.Context) {
	s.server.GracefulStop()
}
