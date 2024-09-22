package config

import (
	"github.com/axidex/Unknown/pkg/db"
	"github.com/axidex/Unknown/pkg/logger"
	"github.com/spf13/viper"
)

type Config struct {
	Server      Server              `yaml:"server"`
	Workdir     string              `yaml:"workdir"`
	Logger      logger.ConfigLogger `yaml:"logger"`
	Postgres    db.Postgres         `yaml:"postgres"`
	Instruments Instruments         `yaml:"instruments"`
	Archive     Archive             `yaml:"archive"`
}

type Archive struct {
	MaxSize    int                 `yaml:"maxSize"`
	Extensions map[string][]string `yaml:"extensions"`
}

type Server struct {
	Port      int       `yaml:"port"`
	Deadlines Deadlines `yaml:"deadlines"`
}

type Deadlines struct {
	SS int `yaml:"ss"`
}

type Instruments struct {
	GitLeaks ShellCommand `yaml:"gitLeaks"`
}

type ShellCommand struct {
	Binary             string   `yaml:"binary"`
	AdditionalCommands []string `yaml:"additionalCommands"`
	Timeout            int      `yaml:"timeout"`
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
