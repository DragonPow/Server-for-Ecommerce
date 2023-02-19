package main

import (
	"Server-for-Ecommerce/app_v2/product_service/cache/mem_cache"
	"Server-for-Ecommerce/app_v2/product_service/cache/redis_cache"
	"Server-for-Ecommerce/app_v2/product_service/config"
	"Server-for-Ecommerce/app_v2/product_service/database/store"
	"Server-for-Ecommerce/app_v2/product_service/internal/service"
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

	if cfg.EnableCache && cfg.EnableRedis {
		go func() {
			err := serviceInstance.Consume()
			if err != nil {
				logger.Error(err, "Consume error")
			}
		}()
	}

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

	var (
		redis    redis_cache.RedisCache
		memCache mem_cache.MemCache
	)

	// Redis cache
	if cfg.EnableCache && cfg.EnableRedis {
		redis = redis_cache.NewCache(
			cfg.RedisConfig.Addr,
			cfg.RedisConfig.Password,
			cfg.RedisConfig.ExpiredDefault,
		)
	} else {
		redis = &redis_cache.NullRedisCache{}
	}

	// Memory cache
	if cfg.EnableCache && cfg.EnableMem {
		//memCache := mem_cache.NewCache(cfg.MemCacheConfig.MaxTimeMiss, cfg.MemCacheConfig.MaxNumberCache)
		memCache, err = mem_cache.NewBigCache(cfg.MemCacheConfig)
		if err != nil {
			logger.Error(err, "Error create mem cache")
			return nil, err
		}
	} else {
		memCache = &mem_cache.NullMemCache{}
	}

	return service.NewService(cfg, logger, store, redis, memCache), nil
}

func newDB(dsn string) (*sqlx.DB, error) {
	db, err := sqlx.Open("postgres", dsn)
	db.SetConnMaxIdleTime(0)
	if err != nil {
		return nil, err
	}
	err = db.Ping()
	if err != nil {
		return nil, err
	}
	return db, nil
}
