package api

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func (app *App) health(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "OK"})
}
