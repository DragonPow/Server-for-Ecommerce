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
	Server          config.ServerConfig `json:"server" mapstructure:"server"`
	MigrationFolder string              `json:"migration_folder" mapstructure:"migration_folder"`
	OrderServiceDB  database.DBConfig   `json:"order_service_db" mapstructure:"order_service_db"`
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
		Server:          config.DefaultServerConfig(),
		MigrationFolder: "file://app/order_service/sql/migrations",
		OrderServiceDB:  database.PostgresSQLDefaultConfig(),
	}
}
