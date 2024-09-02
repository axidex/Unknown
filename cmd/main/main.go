package main

import (
	"Unknown/config"
	"Unknown/internal/api"
	"Unknown/internal/repository"
	"Unknown/pkg/db/tables"
	"Unknown/pkg/logger"
	"fmt"
	"github.com/joho/godotenv"
	"os"
	"strconv"
)

func main() {
	appConfig, err := config.ReadConfig()
	if err != nil {
		fmt.Printf("Got error when reading config from file - %s", err)
		return
	}
	fmt.Printf("Config: %+v\n", appConfig)

	appLogger, err := logger.CreateNewZeroLogger(appConfig.Logger)
	if err != nil {
		fmt.Printf("Got error when initializing logger - %s", err)
		return
	}

	err = godotenv.Load()
	if err != nil {
		appLogger.Warnf("Error loading .env file")

		appConfig.Postgres.User = os.Getenv("DB_USERNAME")
		appConfig.Postgres.Pass = os.Getenv("DB_PASSWORD")
		appConfig.Postgres.Database = os.Getenv("DB_NAME")
		appConfig.Postgres.Port, _ = strconv.Atoi(os.Getenv("DB_PORT"))
		appConfig.Postgres.Url = os.Getenv("DB_URL")

	}

	// Postgres

	appLogger.Infof("Postgres connection - %s:%d/%s/%s", appConfig.Postgres.Url, appConfig.Postgres.Port, appConfig.Postgres.Database, appConfig.Postgres.Schema)

	appRepository, err := repository.CreateNewRepository(appConfig.Postgres)
	if err != nil {
		appLogger.Fatalf("Got error when initializing repository - %s", err)
		return
	}
	appLogger.Info("Creating Schema")
	err = appRepository.CreateSchema(appConfig.Postgres.Schema)
	if err != nil {
		appLogger.Fatalf("Got error when initializing db schema - %s", err)
		return
	}
	appLogger.Info("Running migrations")
	err = appRepository.Migrate(tables.Client{}, tables.Task{})
	if err != nil {
		appLogger.Fatalf("Got error while performing db migrations - %s", err)
		return
	}
	appLogger.Info("Database ready")

	// App
	app := api.CreateApp(appConfig, appLogger)
	engine := app.InitRoutes()
	err = engine.Run(fmt.Sprintf(":%d", appConfig.Server.Port))
	if err != nil {
		appLogger.Fatal("Failed to start server - %s", err)
		return
	}

}
