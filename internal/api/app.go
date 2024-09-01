package api

import (
	"Unknown/config"
	"Unknown/pkg/logger"
	ginzerolog "github.com/dn365/gin-zerolog"
	"github.com/gin-gonic/gin"
)

type App struct {
	config *config.Config
	logger logger.Logger
}

func CreateApp(config *config.Config, logger logger.Logger) *App {
	return &App{
		config: config,
		logger: logger,
	}
}

func (app *App) InitRoutes() *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	router := gin.New()
	router.Use(gin.Recovery())
	router.Use(gin.LoggerWithFormatter(app.loggerMiddleware))

	//router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	api := router.Group("/api")
	{
		health := api.Group("/health")
		{
			health.GET("/ping", app.health)
		}
	}

	router.Use(ginzerolog.Logger("gin"))

	return router
}
