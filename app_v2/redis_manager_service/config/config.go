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
}

type RedisConfig struct {
	Addr                           string `json:"addr" mapstructure:"addr"`
	Password                       string `json:"password" mapstructure:"password"`
	ExpiredDefault                 uint32 `json:"expired_default" mapstructure:"expired_default"`
	MaxRingNumber                  int    `json:"max_ring_number" mapstructure:"max_ring_number"`
	TimeoutRingWriterInMillisecond int    `json:"timeout_ring_writer_in_millisecond" mapstructure:"timeout_ring_writer_in_millisecond"`
}

type KafkaConfig struct {
	UpdateDbConsumer             Consumer `json:"update_db_consumer" mapstructure:"update_db_consumer"`
	Connections                  []string `json:"connections" mapstructure:"connections"`
	MaxPublishTimeoutSecond      int      `json:"max_publish_timeout_second" mapstructure:"max_publish_timeout_second"`
	MaxNumberRetry               int      `json:"max_number_retry" mapstructure:"max_number_retry"`
	TimeSleepPerRetryMillisecond int      `json:"time_sleep_per_retry_millisecond" mapstructure:"time_sleep_per_retry_millisecond"`
}

type Consumer struct {
	Topic       string   `json:"topic" mapstructure:"topic"`
	Connections []string `json:"connections" mapstructure:"connections"`
	Group       string   `json:"group" mapstructure:"group"`
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
		MigrationFolder:  "file://app_v2/redis_manager_service/internal/database/migrations",
		ProductServiceDB: database.PostgresSQLDefaultConfig(),
		RedisConfig: RedisConfig{
			ExpiredDefault:                 180,
			MaxRingNumber:                  10,
			TimeoutRingWriterInMillisecond: 3,
		},
		KafkaConfig: KafkaConfig{
			UpdateDbConsumer: Consumer{
				Topic: "update_product",
				Group: "update_product_consumer",
			},
			MaxPublishTimeoutSecond:      10,
			MaxNumberRetry:               3,
			TimeSleepPerRetryMillisecond: 200,
		},
	}
}
