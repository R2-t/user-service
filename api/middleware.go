package api

import (
	"errors"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"santiagotorres.me/user-service/repositories"
	"santiagotorres.me/user-service/services"
)

func (appState *AppState) CheckJWT() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")

		if authHeader == "" {
			c.AbortWithStatusJSON(401, gin.H{"error": "missing token"})
			return
		}

		const BearerPrefix = "Bearer "

		if !strings.HasPrefix(authHeader, BearerPrefix) {
			c.AbortWithStatusJSON(401, gin.H{"error": "invalid token format"})
			return
		}

		tokenString := strings.TrimPrefix(authHeader, BearerPrefix)

		// Use valid secret key
		token, claims, err := services.ValidateToken(tokenString, "secret")

		if err != nil {
			if errors.Is(err, jwt.ErrTokenExpired) {
				c.AbortWithStatusJSON(401, gin.H{"error": "token expired"})
				return
			}
			c.AbortWithStatusJSON(401, gin.H{"error": "invalid token"})
		}

		userSession, err := repositories.GetUserSession(token.Raw, claims, appState.Db)

		if err != nil {
			c.AbortWithStatusJSON(401, gin.H{"error": "invalid token"})
			return
		}

		if userSession.IsRevoked {
			c.AbortWithStatusJSON(401, gin.H{"error": "token revoked"})
			return
		}

		c.Set("claims", claims)
		c.Set("userSession", userSession)

		c.Next()
	}
}
