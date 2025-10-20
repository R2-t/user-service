package repositories

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"santiagotorres.me/user-service/logger"
	"santiagotorres.me/user-service/models"
	"santiagotorres.me/user-service/services"
	"santiagotorres.me/user-service/utils"
)

func CreateUser(user *models.User, db *gorm.DB) (*uuid.UUID, error) {
	ctx := context.Background()
	_, err := gorm.G[models.User](db).Where("email = ?", user.Email).First(ctx)

	if err == nil {
		return nil, ErrUserAlreadyExists
	}

	if !errors.Is(err, gorm.ErrRecordNotFound) {
		logger.Logger.Error("Error checking user existence", "err", err.Error())
		return nil, NewRepositoryError(ErrCodeDatabaseError, "failed to check user existence", err)
	}

	hashedPassword, err := utils.HashPassword(user.Password)

	if err != nil {
		logger.Logger.Error("Error hashing password", "err", err.Error())
		return nil, NewRepositoryError(ErrCodeHashingError, "failed to hash password", err)
	}

	user.Password = hashedPassword

	createErr := gorm.G[models.User](db).Create(ctx, user)

	if createErr != nil {
		logger.Logger.Error("Error creating user", "err", createErr.Error())
		return nil, NewRepositoryError(ErrCodeDatabaseError, "failed to create user", createErr)
	}

	return &user.UserID, nil
}

// LoginUser logs in a user by email and password.
func LoginUser(email string, userPassword string, deviceInfo *models.DeviceInfo, db *gorm.DB) (*models.PairToken, error) {
	ctx := context.Background()

	user, err := gorm.G[models.User](db).Where("email = ?", email).First(ctx)

	if errors.Is(err, gorm.ErrRecordNotFound) {
		logger.Logger.Warn("Login attempt for non-existent user", "email", email)
		return nil, ErrInvalidCredentials
	}

	if err != nil {
		logger.Logger.Error("Error finding user", "err", err.Error())
		return nil, NewRepositoryError(ErrCodeDatabaseError, "failed to find user", err)
	}

	if !utils.CheckPasswordHash(userPassword, user.Password) {
		logger.Logger.Warn("Invalid password attempt", "email", email)
		return nil, ErrInvalidCredentials
	}

	if !user.UserTOTP.IsEnabled {
		// replace for actual secret
		token, err := services.GenerateTokenWithSession(&user, deviceInfo, db, "secret")
		if err != nil {
			logger.Logger.Error("Error generating token", "err", err.Error())
			return nil, NewRepositoryError(ErrCodeTokenGenerationError, "failed to generate access token", err)
		}
		return token, nil
	}

	// replace for actual secret
	token, err := services.GenerateTempToken(&user, time.Minute*10, "secret")

	if err != nil {
		logger.Logger.Error("Error generating temp token", "err", err.Error())
		return nil, NewRepositoryError(ErrCodeTokenGenerationError, "failed to generate temp auth token", err)
	}

	return &models.PairToken{
		AccessToken:  token,
		RefreshToken: "",
	}, nil
}

func ChangeUserEmail() {}

func ChangePassword() {}

func ForgotPassword() {}

func RefreshToken() {}

func DeleteUser() {}

func LogOutSession() {}

func VerifyJWT() {}

func GenerateTOTP() {}

func VerifyTOTP() {}

func SetUpTOTP() {}
