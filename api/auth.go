package api

import "github.com/gin-gonic/gin"

func (appState *AppState) SetUpAuthRoutes(r *gin.Engine) {
	authRouter := r.Group("/auth")

	{
		authRouter.POST("/register")
		authRouter.POST("/login")
		authRouter.POST("/refresh")
		authRouter.POST("/logout")
		authRouter.POST("/forgot-password")
		authRouter.POST("/reset-password")
		authRouter.POST("/verify")
	}
}
