package database

import (
	"fmt"
	"net/url"
)

type DBConfig interface {
	String() string
	DSN() string
}

type Config struct {
	Host     string `json:"host" mapstructure:"host" yaml:"host"`
	Database string `json:"database" mapstructure:"database" yaml:"database"`
	Port     int    `json:"port" mapstructure:"port" yaml:"port"`
	Username string `json:"username" mapstructure:"username" yaml:"username"`
	Password string `json:"password" mapstructure:"password" yaml:"password"`
	Options  string `json:"options" mapstructure:"options" yaml:"options"`
}

func (c Config) DSN() string {
	options := c.Options
	if options != "" {
		if options[0] != '?' {
			options = "?" + options
		}
	}
	return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s%s",
		c.Username,
		url.QueryEscape(c.Password),
		c.Host,
		c.Port,
		c.Database,
		options)
}

func (c PostgreSQLConfig) DSN() string {
	return fmt.Sprintf("%s:%s@%s:%d/%s%s", c.Username, url.QueryEscape(c.Password), c.Host, c.Port, c.Database, c.Options)
}

type PostgreSQLConfig struct {
	Config `mapstructure:",squash"`
}

func (c PostgreSQLConfig) String() string {
	return fmt.Sprintf("postgresql://%s", c.DSN())
}

func PostgresSQLDefaultConfig() PostgreSQLConfig {
	return PostgreSQLConfig{Config{
		Host:     "127.0.0.1",
		Port:     5432,
		Database: "sample",
		Username: "default",
		Password: "secret",
		Options:  "",
	}}
}
