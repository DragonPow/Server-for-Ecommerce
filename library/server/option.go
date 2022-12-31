package server

type Option func(*Config)

func createConfig(opts []Option) *Config {
	c := createDefaultConfig()
	for _, opt := range opts {
		opt(c)
	}
	return c
}

func WithGatewayAddr(host string, port int) Option {
	return func(config *Config) {
		config.Gateway.Addr = Listen{
			Host: host,
			Port: port,
		}
	}
}

func WithServiceServer(servers ...ServiceServer) Option {
	return func(config *Config) {
		config.ServiceServers = append(config.ServiceServers, servers...)
	}
}

func WithGatewayAddrListen(l Listen) Option {
	return func(c *Config) {
		c.Gateway.Addr = l
	}
}

func WithGrpcAddr(host string, port int) Option {
	return func(c *Config) {
		c.Grpc.Addr = Listen{
			Host: host,
			Port: port,
		}
	}
}

func WithGrpcAddrListen(l Listen) Option {
	return func(c *Config) {
		c.Grpc.Addr = l
	}
}

func WithGatewayServerHandler(handlers ...HTTPServerHandler) Option {
	return func(c *Config) {
		c.Gateway.ServerHandlers = append(c.Gateway.ServerHandlers, handlers...)
	}
}