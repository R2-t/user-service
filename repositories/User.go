package repositories

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/pquerna/otp/totp"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
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

	user, err := gorm.G[models.User](db).Where("email = ?", email).Joins(clause.JoinTarget{Association: "UserTOTP"}, nil).First(ctx)

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

func SetUpTOTP(userId uuid.UUID, userEmail string, issuer string, encryptor *utils.EncryptorManager, db *gorm.DB) (string, error) {
	totpKey, err := totp.Generate(totp.GenerateOpts{
		Issuer:      issuer,
		AccountName: userEmail,
	})

	if err != nil {
		logger.Logger.Error("Error generating TOTP key", "err", err.Error())
		return "", NewRepositoryError(ErrCodeTOTPGenerationError, "failed to generate TOTP key", err)
	}

	totpSecret, err := encryptor.EncryptSecret(totpKey.Secret())

	if err != nil {
		logger.Logger.Error("Error encrypting TOTP secret", "err", err.Error())
		return "", NewRepositoryError(ErrCodeTOTPGenerationError, "failed to encrypt TOTP secret", err)
	}

	totp := models.UserTOTP{
		UserID:    userId,
		Secret:    totpSecret,
		IsEnabled: true,
	}

	ctx := context.Background()
	if err := gorm.G[models.UserTOTP](db).Create(ctx, &totp); err != nil {
		logger.Logger.Error("Error saving User TOTP record", "err", err.Error())
		return "", NewRepositoryError(ErrCodeTOTPGenerationError, "failed to save User TOTP record", err)
	}

	return totpKey.String(), nil
}

func GetUserSession(rawTokenString string, claims *models.Claims, db *gorm.DB) (*models.UserSessions, error) {
	if claims.ExpiresAt != nil && claims.ExpiresAt.Before(time.Now()) {
		logger.Logger.Warn("Expired token received")
		return nil, NewRepositoryError(ErrCodeTokenExpired, "token expired", nil)
	}

	tokenStringHash := utils.HashSHA256(rawTokenString)

	ctx := context.Background()
	userSession, err := gorm.G[models.UserSessions](db).Where("user_id = ? AND token_id = ? AND token_hash = ?", claims.UserID, claims.ID, tokenStringHash).First(ctx)

	if err != nil {
		logger.Logger.Error("An error occur while checking user session", "error", err.Error())
		return nil, NewRepositoryError(ErrCodeInvalidToken, "invalid token", err)
	}

	return &userSession, nil
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
