package main

import (
	"github.com/DragonPow/Server-for-Ecommerce/app/warehouse_service/config"
	"github.com/DragonPow/Server-for-Ecommerce/app/warehouse_service/internal/service"
	"github.com/DragonPow/Server-for-Ecommerce/library/server"
	"github.com/urfave/cli/v2"
	"log"
	"os"
)

var (
	cfg *config.Config
)

func main() {
	if err := run(os.Args); err != nil {
		log.Fatal(err)
	}
}

func run(args []string) (err error) {
	cfg, err = config.Load()
	if err != nil {
		return err
	}
	app := cli.NewApp()
	app.Name = "service"
	app.Commands = []*cli.Command{
		{
			Name:   "server",
			Usage:  "Start grpc/http server",
			Action: serverAction,
		},
	}
	if app.Run(os.Args) != nil {
		panic(err)
	}
	return err
}

func serverAction(context *cli.Context) error {
	service, err := newService(cfg)
	if err != nil {
		log.Printf("Cannot init server, err = %v", err)
		return err
	}
	s, err := server.New(
		server.WithGatewayAddrListen(cfg.Server.HTTP),
		server.WithGrpcAddrListen(cfg.Server.GRPC),
		server.WithServiceServer(service),
	)
	if err != nil {
		log.Printf("Error new server, err = %v", err)
		return err
	}

	if err := s.Serve(); err != nil {
		log.Printf("Error start server, err = %v", err)
		return err
	}
	return nil
}

func newService(cfg *config.Config) (*service.Service, error) {
	return service.NewService(cfg), nil
}
