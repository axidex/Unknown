package repository

import (
	"Unknown/config"
	"Unknown/pkg/db"
	"Unknown/pkg/db/tables"
	"errors"
	"fmt"
	"gorm.io/gorm"
)

type Repository interface {
	Migrate(models ...interface{}) error
	CreateSchema(name string) error
}

type UnknownRepository struct {
	db *gorm.DB
}

func CreateNewRepository(config config.Postgres) (Repository, error) {
	dbConnection, err := db.CreatePostgresConnection(fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=disable", config.Url, config.User, config.Pass, config.Database, config.Port))
	if err != nil {
		return nil, err
	}

	return &UnknownRepository{
		db: dbConnection,
	}, nil
}

func (r *UnknownRepository) Migrate(models ...interface{}) error {
	if models == nil {
		return errors.New("nothing to migrate")
	}
	err := r.db.Migrator().AutoMigrate(&tables.Client{}, &tables.Task{})
	if err != nil {
		return err
	}

	return nil
}

func (r *UnknownRepository) CreateSchema(name string) error {
	sql := fmt.Sprintf(`CREATE SCHEMA IF NOT EXISTS "%s"`, name)
	err := r.db.Exec(sql).Error
	if err != nil {
		return err
	}
	return nil
}
