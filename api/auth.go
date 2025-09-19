package api

import "github.com/gin-gonic/gin"

func (appState *AppState) SetUpAuthRoutes(r *gin.Engine) {
	authRouter := r.Group("/auth")

	{
		authRouter.POST("/register", func(context *gin.Context) {

		})
		authRouter.POST("/login", func(context *gin.Context) {

		})
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
