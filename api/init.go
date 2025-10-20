package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type AppState struct {
	Db *gorm.DB
}

func (appState *AppState) SetupRoutes(r *gin.Engine) {
	r.GET("/health", HealthCheck)
}

func HealthCheck(c *gin.Context) {
	c.String(http.StatusOK, "pong")
}
