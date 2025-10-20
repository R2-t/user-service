package repositories

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"santiagotorres.me/user-service/logger"
	"santiagotorres.me/user-service/models"
	"santiagotorres.me/user-service/utils"
)

func CreateUser(user *models.User, db *gorm.DB) (*uuid.UUID, error) {
	ctx := context.Background()
	_, err := gorm.G[models.User](db).Where("email = ?", user.Email).First(ctx)

	if err == nil {
		return nil, errors.New("user already exists")
	}

	if !errors.Is(err, gorm.ErrRecordNotFound) {
		logger.Logger.Error("Error checking user existence", "err", err.Error())
		return nil, err
	}

	hashedPassword, err := utils.HashPassword(user.Password)

	if err != nil {
		logger.Logger.Error("Error hashing password", "err", err.Error())
		return nil, err
	}

	user.Password = hashedPassword

	createErr := gorm.G[models.User](db).Create(ctx, user)

	if createErr != nil {
		logger.Logger.Error("Error creating user", "err", err.Error())
		return nil, createErr
	}

	return &user.UserID, nil
}

// LoginUser logs in a user by email and password.
func LoginUser(email string, userPassword string, deviceInfo *models.DeviceInfo, db *gorm.DB) (string, error) {
	ctx := context.Background()

	user, err := gorm.G[models.User](db).Where("email = ?", email).First(ctx)

	if errors.Is(err, gorm.ErrRecordNotFound) {
		logger.Logger.Warn("User not found", "err", err.Error())
		return "", err
	}

	if err != nil {
		logger.Logger.Error("Error finding user", "err", err.Error())
		return "", err
	}

	if !utils.CheckPasswordHash(userPassword, user.Password) {
		logger.Logger.Error("Invalid password")
		return "", errors.New("Invalid password")
	}

	if !user.UserTOTP.IsEnabled {
		token, err := utils.GenerateTokenWithSession(&user, models.AccessToken, time.Hour, deviceInfo, db)
		if err != nil {
			logger.Logger.Error("Error generating token", "err", err.Error())
			return "", err
		}
		return token, nil
	}

	return utils.GenerateTempToken(
		&user,
		time.Minute*10,
	)
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
