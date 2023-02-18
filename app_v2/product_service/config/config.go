package config

import (
	"Server-for-Ecommerce/library/config"
	"Server-for-Ecommerce/library/database"
	"bytes"
	"encoding/json"
	"github.com/spf13/viper"
	"log"
	"strings"
)

type Config struct {
	Server           config.ServerConfig `json:"server" mapstructure:"server"`
	MigrationFolder  string              `json:"migration_folder" mapstructure:"migration_folder"`
	ProductServiceDB database.DBConfig   `json:"product_service_db" mapstructure:"product_service_db"`
	RedisConfig      RedisConfig         `json:"redis_config" mapstructure:"redis_config"`
	KafkaConfig      KafkaConfig         `json:"kafka_config" mapstructure:"kafka_config"`
	MemCacheConfig   MemConfig           `json:"mem_cache_config" mapstructure:"mem_cache_config"`
	EnableMem        bool                `json:"enable_mem" mapstructure:"enable_mem"`
	EnableRedis      bool                `json:"enable_redis" mapstructure:"enable_redis"`
	EnableCache      bool                `json:"enable_cache" mapstructure:"enable_cache"`
}

type RedisConfig struct {
	Addr                    string `json:"addr" mapstructure:"addr"`
	Password                string `json:"password" mapstructure:"password"`
	ExpiredDefault          uint32 `json:"expired_default" mapstructure:"expired_default"`
	NumberCachePage         int64  `json:"number_cache_page" mapstructure:"number_cache_page"`
	ExpireCachePageInSecond int    `json:"expire_cache_page_in_second" mapstructure:"expire_cache_page_in_second"`
}

type KafkaConfig struct {
	UpdateDbConsumer Consumer `json:"update_db_consumer" mapstructure:"update_db_consumer"`
	Connections      []string `json:"connections" mapstructure:"connections"`
}

type MemConfig struct {
	MaxTimeMiss                     int `json:"max_time_miss" mapstructure:"max_time_miss"`
	MaxCacheSizeInMB                int `json:"max_cache_size_in_mb" mapstructure:"max_cache_size_in_mb"`
	ExpiredTimeInSecond             int `json:"expired_time_in_second" mapstructure:"expired_time_in_second"`
	TimeBetweenCleanExpiredInSecond int `json:"time_between_clean_expired_in_second" mapstructure:"time_between_clean_expired_in_second"`
	Shards                          int `json:"shards" mapstructure:"shards"`
}

type Consumer struct {
	Topic string `json:"topic" mapstructure:"topic"`
	Group string `json:"group" mapstructure:"group"`
}

// Load system env config
func Load() (*Config, error) {
	/**
	|-------------------------------------------------------------------------
	| hacking to load reflect structure config into env
	|-----------------------------------------------------------------------*/viper.SetConfigType("yaml")
	// You should set default config value here
	c := loadDefaultConfig()

	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.SetConfigName("config")

	err := viper.ReadInConfig()
	if err != nil {
		log.Println("Read config file failed. ", err)

		configBuffer, err := json.Marshal(c)

		if err != nil {
			return nil, err
		}

		viper.ReadConfig(bytes.NewBuffer(configBuffer))
	}

	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "__"))

	// -- end of hacking --//
	viper.AutomaticEnv()
	err = viper.Unmarshal(c)
	return c, err
}

func loadDefaultConfig() *Config {
	return &Config{
		Server:           config.DefaultServerConfig(),
		MigrationFolder:  "file://app_v2/product_service/database/migrations",
		ProductServiceDB: database.PostgresSQLDefaultConfig(),
		RedisConfig: RedisConfig{
			ExpiredDefault:          180,
			NumberCachePage:         5,
			ExpireCachePageInSecond: 60 * 10,
		},
		KafkaConfig: KafkaConfig{
			UpdateDbConsumer: Consumer{
				Topic: "update_cache",
				Group: "product_consume_update_product_consumer",
			},
		},
		MemCacheConfig: MemConfig{
			MaxTimeMiss:                     3,
			MaxCacheSizeInMB:                1024 * 3, // 3 GB
			ExpiredTimeInSecond:             60 * 10,  // 10 minutes
			TimeBetweenCleanExpiredInSecond: 60 * 5,   // 5 minutes
			Shards:                          1024,
		},
	}
}
