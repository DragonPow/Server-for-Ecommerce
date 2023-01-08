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
}

type RedisConfig struct {
	Addr           string `json:"addr" mapstructure:"addr"`
	Password       string `json:"password" mapstructure:"password"`
	ExpiredDefault uint32 `json:"expired_default" mapstructure:"expired_default"`
}

type KafkaConfig struct {
	UpdateDbConsumer Consumer `json:"consumers" mapstructure:"consumers"`
}

type MemConfig struct {
	MaxTimeMiss int `json:"max_time_miss" mapstructure:"max_time_miss"`
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
		MigrationFolder:  "file://app_v2/product_service/internal/database/migrations",
		ProductServiceDB: database.PostgresSQLDefaultConfig(),
	}
}