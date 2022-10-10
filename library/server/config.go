package server

import (
	"fmt"
	"net"
)

type Config struct {
	Gateway        *gatewayConfig
	Grpc           *grpcConfig
	ServiceServers []ServiceServer
}

type Listen struct {
	Host string `json:"host" yaml:"host"`
	Port int    `json:"port" yaml:"port"`
}

func (l Listen) String() string {
	return fmt.Sprintf("%s:%v", l.Host, l.Port)
}

func (l Listen) CreateListener() (net.Listener, error) {
	lis, err := net.Listen("tcp", l.String())
	if err != nil {
		return nil, fmt.Errorf("Fail to listen, %v", l.String())
	}
	return lis, err
}

func createDefaultConfig() *Config {
	return &Config{
		Gateway:        createDefaultGatewayConfig(),
		Grpc:           createDefaultGrpcConfig(),
		ServiceServers: nil,
	}
}
