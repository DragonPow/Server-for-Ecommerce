package main

import (
	"context"
	accountApi "github.com/DragonPow/Server-for-Ecommerce/app/account_service/api"
	orderApi "github.com/DragonPow/Server-for-Ecommerce/app/order_service/api"
	"github.com/DragonPow/Server-for-Ecommerce/app/warehouse_service/config"
	"github.com/DragonPow/Server-for-Ecommerce/app/warehouse_service/internal/service"
	"github.com/DragonPow/Server-for-Ecommerce/app/warehouse_service/internal/store"
	"github.com/DragonPow/Server-for-Ecommerce/library/database/migrate"
	"github.com/DragonPow/Server-for-Ecommerce/library/log"
	"github.com/DragonPow/Server-for-Ecommerce/library/server"
	"github.com/go-logr/logr"
	"github.com/jmoiron/sqlx"
	"github.com/urfave/cli/v2"
	"google.golang.org/grpc"
	"os"
)

var (
	cfg    *config.Config
	logger logr.Logger
)

func main() {
	logger = log.MustBuildLogR()
	if err := run(os.Args); err != nil {
		logger.Error(err, "Run server fail")
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
		{
			Name:        "migrate",
			Usage:       "doing database migration",
			Subcommands: migrate.CliCommand(cfg.MigrationFolder, cfg.WarehouseServiceDB.String()),
		},
	}
	if app.Run(args) != nil {
		panic(err)
	}
	return err
}

func serverAction(context *cli.Context) error {
	service, err := newService(cfg)
	if err != nil {
		logger.Error(err, "Cannot init server")
		return err
	}
	s, err := server.New(
		server.WithGatewayAddrListen(cfg.Server.HTTP),
		server.WithGrpcAddrListen(cfg.Server.GRPC),
		server.WithServiceServer(service),
	)
	if err != nil {
		logger.Error(err, "Error new server")
		return err
	}

	if err := s.Serve(); err != nil {
		logger.Error(err, "Error start server")
		return err
	}
	return nil
}

func newService(cfg *config.Config) (*service.Service, error) {
	db, err := newDB(cfg.WarehouseServiceDB.String())
	if err != nil {
		logger.Error(err, "Error connect database")
		return nil, err
	}
	store := store.NewStore(db, logger)

	// Order Client
	orderClientConnect, err := grpc.DialContext(context.Background(), cfg.OrderServiceAddr, grpc.WithInsecure())
	orderClient := orderApi.NewOrderServiceClient(orderClientConnect)

	// AccountClient
	accountClientConnect, err := grpc.DialContext(context.Background(), cfg.AccountServiceAddr, grpc.WithInsecure())
	accountClient := accountApi.NewAccountServiceClient(accountClientConnect)

	return service.NewService(cfg, logger, store, orderClient, accountClient), nil
}

func newDB(dsn string) (*sqlx.DB, error) {
	db, err := sqlx.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}
	err = db.Ping()
	if err != nil {
		return nil, err
	}
	return db, nil
}
