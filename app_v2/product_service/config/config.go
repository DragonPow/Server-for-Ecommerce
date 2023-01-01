package config

import (
	"bytes"
	"encoding/json"
	"github.com/DragonPow/Server-for-Ecommerce/library/config"
	"github.com/DragonPow/Server-for-Ecommerce/library/database"
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
	Addr           string `json:"addr" mapstructure:"addr"`
	Password       string `json:"password" mapstructure:"password"`
	ExpiredDefault uint32 `json:"expired_default" mapstructure:"expired_default"`
}

type KafkaConfig struct {
	Topic         string   `json:"topic" mapstructure:"topic"`
	Connections   []string `json:"connections" mapstructure:"connections"`
	ConsumerGroup string   `json:"consumer_group" mapstructure:"consumer_group"`
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
		MigrationFolder:  "file://app_v2/product_service/internal/database/migrations",
		ProductServiceDB: database.PostgresSQLDefaultConfig(),
		RedisConfig: RedisConfig{
			Addr:           "localhost:6379",
			Password:       "",
			ExpiredDefault: 0,
		},
	}
}
