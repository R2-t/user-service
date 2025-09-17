package api

import (
	"errors"
	"net/http"
	"santiagotorres.me/user-service/utils"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"santiagotorres.me/user-service/logger"
	"santiagotorres.me/user-service/models"
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

	_, err := gorm.G[models.User](appState.Db).Where("email = ?", user.Email).First(c.Request.Context())

	if err == nil {
		c.JSON(http.StatusConflict, gin.H{
			"error": "User with this email already exists",
		})
		return
	}

	if !errors.Is(err, gorm.ErrRecordNotFound) {
		logger.Logger.ErrorContext(c.Request.Context(), "Error checking user existence", "err", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{})
		return
	}

	hashedPassword, err := utils.HashPassword(user.Password)

	if err != nil {
		logger.Logger.ErrorContext(c.Request.Context(), "Error hashing password", "err", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{})
		return
	}

	user.Password = hashedPassword

	createErr := gorm.G[models.User](appState.Db).Create(c.Request.Context(), &user)

	if createErr != nil {
		logger.Logger.ErrorContext(c.Request.Context(), "Error creating user", "err", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"userId": user.ID})
}

func (appState *AppState) LogIn(c *gin.Context) {

}
