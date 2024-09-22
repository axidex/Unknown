package db

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func CreatePostgresConnection(connectionString string) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(connectionString), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	return db, nil
}

type Postgres struct {
	Url      string `yaml:"url"`
	Port     int    `yaml:"port"`
	User     string `yaml:"user"`
	Pass     string `yaml:"pass"`
	Database string `yaml:"database"`
	Schema   string `yaml:"schema"`
}
