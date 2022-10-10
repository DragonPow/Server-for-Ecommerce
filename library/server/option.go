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
