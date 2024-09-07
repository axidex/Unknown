package config

import (
	"github.com/axidex/Unknown/pkg/db"
	"github.com/axidex/Unknown/pkg/logger"
	"github.com/spf13/viper"
)

type Config struct {
	Server   Server              `yaml:"server"`
	Logger   logger.ConfigLogger `yaml:"logger"`
	Postgres db.Postgres         `yaml:"postgres"`
}

type Server struct {
	Port int `yaml:"port"`
}

func ReadConfig() (*Config, error) {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./")
	viper.AddConfigPath("./config/")

	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return nil, err
	}

	return &config, nil
}
