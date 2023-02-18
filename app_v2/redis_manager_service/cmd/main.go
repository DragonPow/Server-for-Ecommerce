package main

import (
	"Server-for-Ecommerce/app_v2/redis_manager_service/config"
	"Server-for-Ecommerce/app_v2/redis_manager_service/internal/database/store"
	"Server-for-Ecommerce/app_v2/redis_manager_service/internal/redis"
	"Server-for-Ecommerce/app_v2/redis_manager_service/internal/service"
	"Server-for-Ecommerce/app_v2/redis_manager_service/util"
	"Server-for-Ecommerce/library/database/migrate"
	pub "Server-for-Ecommerce/library/kafka/pub"
	"Server-for-Ecommerce/library/log"
	"Server-for-Ecommerce/library/server"
	"github.com/go-logr/logr"
	"github.com/jmoiron/sqlx"
	"github.com/urfave/cli/v2"
	"os"
	"time"
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
	logger.Info("Config success", "cfg", cfg)

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
			Subcommands: migrate.CliCommand(cfg.MigrationFolder, cfg.ProductServiceDB.String()),
		},
	}
	if app.Run(args) != nil {
		panic(err)
	}
	return err
}

func serverAction(context *cli.Context) error {
	serviceInstance, err := newService(cfg)
	if err != nil {
		logger.Error(err, "Cannot init server")
		return err
	}
	s, err := server.New(
		server.WithGatewayAddrListen(cfg.Server.HTTP),
		server.WithGrpcAddrListen(cfg.Server.GRPC),
		server.WithServiceServer(serviceInstance),
	)
	if err != nil {
		logger.Error(err, "Error new server")
		return err
	}

	go func() {
		err := serviceInstance.Consume()
		if err != nil {
			logger.Error(err, "Consume error")
		}
	}()

	if err := s.Serve(); err != nil {
		logger.Error(err, "Error start server")
		return err
	}
	return nil
}

func newService(cfg *config.Config) (*service.Service, error) {
	db, err := newDB(cfg.ProductServiceDB.String())
	if err != nil {
		logger.Error(err, "Error connect database")
		return nil, err
	}
	store := store.NewStore(db, logger)

	redis := redis.New(cfg.RedisConfig)

	// Producer
	producerCfg := cfg.KafkaConfig
	producer, err := pub.NewProducer(
		producerCfg.Connections,
		&logger,
		pub.WithPublishTimeout(time.Duration(producerCfg.MaxPublishTimeoutSecond)*time.Second),
		pub.WithMaxNumberRetry(producerCfg.MaxNumberRetry),
		pub.WithTimeSleepPerRetry(time.Duration(producerCfg.TimeSleepPerRetryMillisecond)*time.Millisecond),
	)
	if err != nil {
		logger.Error(err, "Create producer fail")
		return nil, err
	}
	err = producer.Register(util.TopicUpdateCache)
	if err != nil {
		logger.Error(err, "Register topic producer fail", "topicName", util.TopicUpdateCache)
		return nil, err
	}

	return service.NewService(cfg, logger, store, redis, producer), nil
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
