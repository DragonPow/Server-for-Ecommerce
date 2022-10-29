package server

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type Server struct {
	grpcServer    *grpcServer
	gatewayServer *gatewayServer
	config        *Config
}

func New(opts ...Option) (*Server, error) {
	c := createConfig(opts)

	log.Println("Create grpc server")
	grpcServer := newGrpcServer(c.Grpc, c.ServiceServers)
	reflection.Register(grpcServer.server)

	conn, err := grpc.Dial(c.Grpc.Addr.String(), grpc.WithInsecure())
	if err != nil {
		return nil, fmt.Errorf("Fail dial Grpc server, %v\n", err)
	}

	log.Println("Create gateway server")
	gatewayServer, err := newGatewayServer(c.Gateway, conn, c.ServiceServers)
	if err != nil {
		return nil, fmt.Errorf("Fail dial Gateway server, %v\n", err)
	}

	return &Server{
		grpcServer:    grpcServer,
		gatewayServer: gatewayServer,
		config:        c,
	}, nil
}

func (s *Server) Serve() error {
	stop := make(chan os.Signal, 1)
	errCh := make(chan error)
	signal.Notify(stop, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		if err := s.gatewayServer.Serve(); err != nil {
			log.Printf("Error starting http server, err = %v", err)
			errCh <- err
		}
	}()
	go func() {
		if err := s.grpcServer.Serve(); err != nil {
			log.Printf("Error starting grpc server, err = %v", err)
			errCh <- err
		}
	}()

	for {
		select {
		case <-stop:
			log.Println("Shutting down server")
			ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
			defer cancel()
			for _, ss := range s.config.ServiceServers {
				ss.Close(ctx)
			}
			s.gatewayServer.Shutdown(ctx)
			s.grpcServer.Shutdown(ctx)
			return nil
		case err := <-errCh:
			return err
		}

	}
}