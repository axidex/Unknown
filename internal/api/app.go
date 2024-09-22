package api

import (
	"github.com/axidex/Unknown/config"
	_ "github.com/axidex/Unknown/docs"
	"github.com/axidex/Unknown/internal/parser"
	"github.com/axidex/Unknown/pkg/archive"
	"github.com/axidex/Unknown/pkg/logger"
	"github.com/axidex/Unknown/pkg/shell"
	"github.com/gin-gonic/gin"
)

type App struct {
	config          *config.Config
	logger          logger.Logger
	archiveManagers map[string]archive.Manager
	parsers         map[string]parser.Parser
	services        map[string]shell.Service
}

func CreateApp(config *config.Config, logger logger.Logger, archiveManagers map[string]archive.Manager, parsers map[string]parser.Parser, services map[string]shell.Service) *App {
	return &App{
		config:          config,
		logger:          logger,
		archiveManagers: archiveManagers,
		parsers:         parsers,
		services:        services,
	}
}

func (app *App) InitRoutes() *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	router := gin.New()
	router.Use(CustomRecoveryFunc(app.logger))
	router.Use(LoggerMiddleware(app.logger))

	//router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	router.GET("/swagger/*any", app.swagger)

	api := router.Group("/api")
	{
		health := api.Group("/health")
		{
			health.GET("/ping", app.health)
		}

		v1 := api.Group("/v1")
		{
			v1.POST("/scan", app.Scan)
		}

	}

	app.listRoutes(router)
	return router
}

func (app *App) listRoutes(router *gin.Engine) {
	for _, route := range router.Routes() {
		app.logger.Infof("Method: %s | Path: %s | Handler: %s", route.Method, route.Path, route.Handler)
	}
}
