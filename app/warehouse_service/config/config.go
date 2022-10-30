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
	Server             config.ServerConfig `json:"server" mapstructure:"server"`
	MigrationFolder    string              `json:"migration_folder" mapstructure:"migration_folder"`
	WarehouseServiceDB database.DBConfig   `json:"warehouse_service_db" mapstructure:"warehouse_service_db"`
	OrderServiceAddr   string              `json:"order_service_addr" mapstructure:"order_service_addr"`
	AccountServiceAddr string              `json:"account_service_addr" mapstructure:"account_service_addr"`
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
		Server:             config.DefaultServerConfig(),
		MigrationFolder:    "file://app/warehouse_service/sql/migrations",
		WarehouseServiceDB: database.PostgresSQLDefaultConfig(),
		OrderServiceAddr:   "order-service-api:443",
		AccountServiceAddr: "account-service-api:443",
	}
}
