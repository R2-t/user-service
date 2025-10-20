package utils

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"time"

	"crypto/sha256"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/datatypes"
	"gorm.io/gorm"
	"santiagotorres.me/user-service/logger"
	"santiagotorres.me/user-service/models"
)

func HashPassword(password string) (string, error) {
	var passwordBytes = []byte(password)

	hashedPasswordBytes, err := bcrypt.GenerateFromPassword(passwordBytes, bcrypt.DefaultCost)

	return string(hashedPasswordBytes), err
}

func CheckPasswordHash(password string, hashedPassword string) bool {
	var passwordBytes = []byte(password)
	var hashedPasswordBytes = []byte(hashedPassword)

	err := bcrypt.CompareHashAndPassword(hashedPasswordBytes, passwordBytes)

	return err == nil
}

func GenerateTempToken(user *models.User, duration time.Duration) (string, error) {
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

	// Replace for actual secret key
	tokenString, err := token.SignedString([]byte("secret"))
	if err != nil {
		logger.Logger.Error("Failed to sign token", "error", err)
		return "", err
	}
	return tokenString, nil
}

func GenerateTokenWithSession(user *models.User, duration time.Duration, deviceInfo *models.DeviceInfo, db *gorm.DB) (string, error) {
	ctx := context.Background()
	jti := uuid.New().String()
	now := time.Now().UTC()
	expiresAt := now.Add(duration)

	// Replace issuer with actual service name
	claims := models.Claims{
		UserID:       user.UserID,
		Email:        user.Email,
		TokenType:    models.AccessToken,
		TOTPVerified: true,
		RegisteredClaims: jwt.RegisteredClaims{
			ID:        jti,
			Subject:   user.UserID.String(),
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(expiresAt),
			Issuer:    "",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Replace for actual secret key
	tokenString, err := token.SignedString([]byte("secret"))
	if err != nil {
		logger.Logger.Error("Failed to sign token", "error", err)
		return "", err
	}

	tokenHash := hashToken(tokenString)

	deviceJson, err := json.Marshal(deviceInfo)
	if err != nil {
		logger.Logger.Error("Failed to marshal device info", "error", err)
		return "", err
	}

	userSession := models.UserSessions{
		UserID:     user.UserID,
		TokenHash:  tokenHash,
		DeviceInfo: datatypes.JSON(deviceJson),
	}

	createSessionErr := gorm.G[models.UserSessions](db).Create(ctx, &userSession)
	if createSessionErr != nil {
		logger.Logger.Error("Failed to create user session", "error", createSessionErr)
		return "", err
	}

	return tokenString, nil
}

func hashToken(token string) string {
	hash := sha256.Sum256([]byte(token))
	return hex.EncodeToString(hash[:])
}
