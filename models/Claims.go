package models

import (
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

const (
	TempAuth     = "temp_auth"
	AccessToken  = "access_token"
	RefreshToken = "refresh_token"
)

type Claims struct {
	UserID       uuid.UUID `json:"user_id"`
	Email        string    `json:"email"`
	TokenType    string    `json:"type"` // "temp_auth" or "access_token" or "refresh_token"
	TOTPVerified bool      `json:"totp_verified"`
	jwt.RegisteredClaims
}

type PairToken struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}
