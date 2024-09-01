package config

import (
	"github.com/spf13/viper"
)

type Config struct {
	Server Server `yaml:"server"`
	Logger Logger `yaml:"logger"`
}

type Server struct {
	Port int `yaml:"port"`
}

type Logger struct {
	Level    string `yaml:"level"`
	FilePath string `yaml:"filePath"`
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
