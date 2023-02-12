package main

import (
	"Server-for-Ecommerce/app/account_service/config"
	"Server-for-Ecommerce/app/account_service/internal/service"
	"Server-for-Ecommerce/app/account_service/internal/store"
	"Server-for-Ecommerce/library/database/migrate"
	"Server-for-Ecommerce/library/log"
	"Server-for-Ecommerce/library/server"
	"github.com/go-logr/logr"
	"github.com/jmoiron/sqlx"
	"github.com/urfave/cli/v2"
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
			Subcommands: migrate.CliCommand(cfg.MigrationFolder, cfg.AccountServiceDB.String()),
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
	db, err := newDB(cfg.AccountServiceDB.String())
	if err != nil {
		logger.Error(err, "Error connect database")
		return nil, err
	}
	store := store.NewStore(db, logger)

	// TODO: add another service here and pass to NewService
	// ...

	return service.NewService(cfg, logger, store), nil
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
