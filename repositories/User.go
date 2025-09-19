package repositories

import (
	"context"
	"errors"
	"gorm.io/gorm"
	"santiagotorres.me/user-service/api"
	"santiagotorres.me/user-service/logger"
	"santiagotorres.me/user-service/models"
	"santiagotorres.me/user-service/utils"
)

func CreateUser(user models.User, appState api.AppState) {
	ctx := context.Background()
	_, err := gorm.G[models.User](appState.Db).Where("email = ?", user.Email).First(ctx)

	if err == nil {
		return
	}

	if !errors.Is(err, gorm.ErrRecordNotFound) {
		logger.Logger.Error("Error checking user existence", "err", err.Error())
		return
	}

	hashedPassword, err := utils.HashPassword(user.Password)

	if err != nil {
		logger.Logger.Error("Error hashing password", "err", err.Error())
		return
	}

	user.Password = hashedPassword

	createErr := gorm.G[models.User](appState.Db).Create(ctx, &user)

	if createErr != nil {
		logger.Logger.Error("Error creating user", "err", err.Error())
		return
	}
}

func LoginUser(userName string, userPassword string, appState api.AppState) {}

func ChangeUserEmail() {}

func ChangePassword() {}

func ForgotPassword() {}

func RefreshToken() {}

func DeleteUser() {}

func LogOutSession() {}

func VerifyJWT() {}

func GenerateTOTP() {}

func VerifyTOTP() {}
