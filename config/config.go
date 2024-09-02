package config

import (
	"github.com/spf13/viper"
)

type Config struct {
	Server   Server   `yaml:"server"`
	Logger   Logger   `yaml:"logger"`
	Postgres Postgres `yaml:"postgres"`
}

type Server struct {
	Port int `yaml:"port"`
}

type Logger struct {
	Level    string `yaml:"level"`
	FilePath string `yaml:"filePath"`
}

type Postgres struct {
	Url      string `yaml:"url"`
	Port     int    `yaml:"port"`
	User     string `yaml:"user"`
	Pass     string `yaml:"pass"`
	Database string `yaml:"database"`
	Schema   string `yaml:"schema"`
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
