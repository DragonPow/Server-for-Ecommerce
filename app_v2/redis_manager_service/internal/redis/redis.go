package redis

import (
	"Server-for-Ecommerce/app_v2/redis_manager_service/config"
	"Server-for-Ecommerce/library/cache/redis"
	"Server-for-Ecommerce/library/ring"
	"github.com/segmentio/kafka-go"
)

type Redis struct {
	*redis.Redis
	Ring                      *ring.RingBuffer[kafka.Message]
	MaxRingNumber             int
	TimeoutRingWriterInSecond int
}

func New(cfg config.RedisConfig) *Redis {
	return &Redis{
		Redis:                     redis.New(cfg.Addr, cfg.Password, cfg.ExpiredDefault),
		Ring:                      ring.New[kafka.Message](cfg.MaxRingNumber),
		MaxRingNumber:             cfg.MaxRingNumber,
		TimeoutRingWriterInSecond: cfg.TimeoutRingWriterInSecond,
	}
}
