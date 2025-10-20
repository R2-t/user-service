package services

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"gorm.io/datatypes"
	"gorm.io/gorm"
	"santiagotorres.me/user-service/logger"
	"santiagotorres.me/user-service/models"
	"santiagotorres.me/user-service/utils"
)

// GenerateTempToken generates a temporary authentication token for TOTP verification.
func GenerateTempToken(user *models.User, duration time.Duration, jwtSecret string) (string, error) {
	jti := uuid.New().String()
	now := time.Now().UTC()
	expiresAt := now.Add(duration)

	claim := models.Claims{
		UserID:       user.UserID,
		Email:        user.Email,
		TokenType:    models.TempAuth,
		TOTPVerified: false,
		RegisteredClaims: jwt.RegisteredClaims{
			ID:        jti,
			Subject:   user.UserID.String(),
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(expiresAt),
			Issuer:    "",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)
	tokenString, err := token.SignedString([]byte(jwtSecret))
	if err != nil {
		logger.Logger.Error("Failed to sign temp token", "error", err)
		return "", err
	}

	return tokenString, nil
}

// GenerateTokenWithSession generates access and refresh tokens and creates a user session.
func GenerateTokenWithSession(
	user *models.User,
	deviceInfo *models.DeviceInfo,
	db *gorm.DB,
	jwtSecret string,
) (*models.PairToken, error) {
	ctx := context.Background()
	tokenJti := uuid.New().String()
	refreshTokenJti := uuid.New().String()
	now := time.Now().UTC()
	tokenExpiresAt := now.Add(time.Hour * 24)
	refreshTokenExpiresAt := now.Add(time.Hour * 7 * 24)

	// Create access token claims
	accessClaims := models.Claims{
		UserID:       user.UserID,
		Email:        user.Email,
		TokenType:    models.AccessToken,
		TOTPVerified: true,
		RegisteredClaims: jwt.RegisteredClaims{
			ID:        tokenJti,
			Subject:   user.UserID.String(),
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(tokenExpiresAt),
			Issuer:    "",
		},
	}

	// Create refresh token claims
	refreshClaims := models.Claims{
		UserID:       user.UserID,
		Email:        user.Email,
		TokenType:    models.RefreshToken,
		TOTPVerified: true,
		RegisteredClaims: jwt.RegisteredClaims{
			ID:        refreshTokenJti,
			Subject:   user.UserID.String(),
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(refreshTokenExpiresAt),
			Issuer:    "",
		},
	}

	// Sign tokens
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims)
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)

	accessTokenString, accessTokenErr := accessToken.SignedString([]byte(jwtSecret))
	refreshTokenString, refreshTokenErr := refreshToken.SignedString([]byte(jwtSecret))

	if accessTokenErr != nil || refreshTokenErr != nil {
		logger.Logger.Error("Failed to sign tokens", "accessTokenError", accessTokenErr, "refreshTokenError", refreshTokenErr)
		return nil, errors.New("failed to sign tokens")
	}

	// Hash tokens for storage
	accessTokenHash := utils.HashSHA256(accessTokenString)
	refreshTokenHash := utils.HashSHA256(refreshTokenString)

	// Marshal device info
	deviceJSON, err := json.Marshal(deviceInfo)
	if err != nil {
		logger.Logger.Error("Failed to marshal device info", "error", err)
		return nil, err
	}

	// Create user session
	userSession := models.UserSessions{
		UserID:           user.UserID,
		TokenID:          tokenJti,
		RefreshTokenID:   refreshTokenJti,
		TokenHash:        accessTokenHash,
		RefreshTokenHash: refreshTokenHash,
		DeviceInfo:       datatypes.JSON(deviceJSON),
	}

	createSessionErr := gorm.G[models.UserSessions](db).Create(ctx, &userSession)
	if createSessionErr != nil {
		logger.Logger.Error("Failed to create user session", "error", createSessionErr)
		return nil, createSessionErr
	}

	return &models.PairToken{
		AccessToken:  accessTokenString,
		RefreshToken: refreshTokenString,
	}, nil
}

// ValidateToken validates a JWT token and returns the claims.
func ValidateToken(tokenString string, jwtSecret string) (*models.Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &models.Claims{}, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(jwtSecret), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*models.Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("invalid token")
}

// RevokeSession revokes a user session by token ID.
//func RevokeSession(tokenID string, db *gorm.DB) error {
//	ctx := context.Background()
//
//	session, err := gorm.G[models.UserSessions](db).Where("token_id = ?", tokenID).First(ctx)
//	if err != nil {
//		logger.Logger.Error("Failed to find session", "error", err)
//		return err
//	}
//
//	deleteErr := gorm.G[models.UserSessions](db).Delete(ctx, session)
//	if deleteErr != nil {
//		logger.Logger.Error("Failed to delete session", "error", deleteErr)
//		return deleteErr
//	}
//
//	return nil
//}
