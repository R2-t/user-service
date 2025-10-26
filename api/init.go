package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"santiagotorres.me/user-service/utils"
)

type AppState struct {
	Db               *gorm.DB
	EncryptorManager *utils.EncryptorManager
}

func (appState *AppState) SetupRoutes(r *gin.Engine) {
	r.GET("/health", HealthCheck)
}

func HealthCheck(c *gin.Context) {
	c.String(http.StatusOK, "pong")
}
