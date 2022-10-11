package server

import (
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"log"
)

type Server struct {
	grpcServer *grpcServer
	gateServer *gatewayServer
	config     *Config
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
		grpcServer: grpcServer,
		gateServer: gatewayServer,
		config:     c,
	}, nil
}
