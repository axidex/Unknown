package main

import (
	"Unknown/config"
	"Unknown/internal/api"
	"Unknown/pkg/logger"
	"fmt"
)

func main() {
	appConfig, err := config.ReadConfig()
	if err != nil {
		fmt.Printf("Got error when reading config from file - %s", err)
		return
	}
	fmt.Printf("Config: %+v\n", appConfig)

	zeroLogger, err := logger.CreateNewZeroLogger(appConfig.Logger)
	if err != nil {
		fmt.Printf("Got error when initializing logger - %s", err)
		return
	}

	app := api.CreateApp(appConfig, zeroLogger)
	engine := app.InitRoutes()
	err = engine.Run(fmt.Sprintf(":%d", appConfig.Server.Port))
	if err != nil {
		zeroLogger.Fatal("Failed to start server - %s", err)
		return
	}

}
