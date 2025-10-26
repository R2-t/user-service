package api

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"santiagotorres.me/user-service/logger"
	"santiagotorres.me/user-service/models"
	"santiagotorres.me/user-service/repositories"
	"santiagotorres.me/user-service/utils"
)

func (appState *AppState) SetUpAuthRoutes(r *gin.Engine) {
	authRouter := r.Group("/auth")

	{
		authRouter.POST("/signup", appState.SignUp)
		authRouter.POST("/login", appState.Login)
		authRouter.POST("/register-totp", appState.CheckJWT(), appState.RegisterTOTP)
		authRouter.POST("/refresh", func(context *gin.Context) {

		})
		authRouter.POST("/logout", func(context *gin.Context) {

		})
		authRouter.POST("/forgot-password", func(context *gin.Context) {

		})
		authRouter.POST("/reset-password", func(context *gin.Context) {

		})
		authRouter.POST("/verify", func(context *gin.Context) {

		})
	}
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

	userId, err := repositories.CreateUser(&user, appState.Db)

	if err != nil {
		// Check for specific error types using errors.Is
		if errors.Is(err, repositories.ErrUserAlreadyExists) {
			c.JSON(http.StatusConflict, gin.H{
				"error": "User with this email already exists",
			})
			return
		}

		// Check for repository errors by type
		var repoErr *repositories.RepositoryError
		if errors.As(err, &repoErr) {
			switch repoErr.Code {
			case repositories.ErrCodeHashingError:
				logger.Logger.ErrorContext(c.Request.Context(), "Password hashing error", "err", err.Error())
				c.JSON(http.StatusInternalServerError, gin.H{
					"error": "Failed to process password",
				})
				return
			case repositories.ErrCodeDatabaseError:
				logger.Logger.ErrorContext(c.Request.Context(), "Database error during signup", "err", err.Error())
				c.JSON(http.StatusInternalServerError, gin.H{
					"error": "An error occurred during registration",
				})
				return
			}
		}

		// Generic error fallback
		logger.Logger.ErrorContext(c.Request.Context(), "Unexpected error during signup", "err", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "An unexpected error occurred",
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"userID": userId})
}

func (appState *AppState) Login(context *gin.Context) {
	var req struct {
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required"`
	}

	if err := context.ShouldBindJSON(&req); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	deviceInfo := utils.ExtractDeviceInfo(context)

	tokens, err := repositories.LoginUser(req.Email, req.Password, &deviceInfo, appState.Db)

	if err != nil {
		// Check for specific error types using errors.Is
		if errors.Is(err, repositories.ErrInvalidCredentials) {
			context.JSON(http.StatusUnauthorized, gin.H{
				"error": "Invalid email or password",
			})
			return
		}

		// Check for repository errors by type
		var repoErr *repositories.RepositoryError
		if errors.As(err, &repoErr) {
			switch repoErr.Code {
			case repositories.ErrCodeDatabaseError:
				logger.Logger.ErrorContext(context.Request.Context(), "Database error during login", "err", err.Error())
				context.JSON(http.StatusInternalServerError, gin.H{
					"error": "An error occurred during login",
				})
				return
			case repositories.ErrCodeTokenGenerationError:
				logger.Logger.ErrorContext(context.Request.Context(), "Token generation error", "err", err.Error())
				context.JSON(http.StatusInternalServerError, gin.H{
					"error": "Failed to generate authentication token",
				})
				return
			}
		}

		// Generic error fallback
		logger.Logger.ErrorContext(context.Request.Context(), "Unexpected error during login", "err", err.Error())
		context.JSON(http.StatusInternalServerError, gin.H{
			"error": "An unexpected error occurred",
		})
		return
	}

	context.JSON(http.StatusOK, tokens)
}

func (appState *AppState) RegisterTOTP(context *gin.Context) {
	// Implement TOTP registration logic here
}
