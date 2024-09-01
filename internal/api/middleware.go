package api

import (
	"github.com/gin-gonic/gin"
)

func (app *App) loggerMiddleware(param gin.LogFormatterParams) string {

	app.logger.Infof("%s \"%s %s %s [%d] %s \"%s\" %s\"",
		param.ClientIP,
		param.Method,
		param.Path,
		param.Request.Proto,
		param.StatusCode,
		param.Latency,
		param.Request.UserAgent(),
		param.ErrorMessage,
	)

	return ""
}
