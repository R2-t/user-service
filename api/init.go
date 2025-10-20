package api

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
	"santiagotorres.me/user-service/logger"
	"santiagotorres.me/user-service/models"
	"santiagotorres.me/user-service/repositories"
)

type AppState struct {
	Db *gorm.DB
}

func (appState *AppState) SetupRoutes(r *gin.Engine) {
	r.GET("/health", HealthCheck)
	r.POST("/signup", appState.SignUp)
}

func HealthCheck(c *gin.Context) {
	c.String(http.StatusOK, "pong")
}

func (appState *AppState) SignUp(c *gin.Context) {
	var user models.User

	if err := c.ShouldBindJSON(&user); err != nil {
		logger.Logger.ErrorContext(
			c.Request.Context(),
			"Error deserializing user",
			"err", err.Error(),
		)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	userId, err := repositories.CreateUser(&user, appState)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{})
	}

	c.JSON(http.StatusCreated, gin.H{"userID": userId})
}

func (appState *AppState) LogIn(c *gin.Context) {

}
