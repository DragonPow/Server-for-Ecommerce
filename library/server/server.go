package server

type Server struct {
	grpcServer *grpcServer
	gateServer *gatewayServer
	config     *Config
}

func New(opts ...Option) (*Server, error) {
	c := createConfig(opts)

	return &Server{
		grpcServer: nil,
		gateServer: nil,
		config:     c,
	}, nil
}
