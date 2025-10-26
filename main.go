package main

import (
	"encoding/base64"
	"fmt"

	"github.com/gin-gonic/gin"
	"santiagotorres.me/user-service/api"
	"santiagotorres.me/user-service/configs"
	"santiagotorres.me/user-service/logger"
	"santiagotorres.me/user-service/utils"
)

func main() {
	settings := configs.GetSettings()

	logger.Logger.Info("Starting service", "service", settings.ServiceName, "port", settings.ServicePort)

	dbConfig := configs.DatabaseConfig{
		Host:     settings.DbHost,
		Port:     settings.DbPort,
		User:     settings.DbUser,
		Password: settings.DbPassword,
		DBName:   settings.DBName,
	}

	decodedSecret, base64Err := base64.StdEncoding.DecodeString(settings.EncryptionKey)

	if base64Err != nil {
		logger.Logger.Error("Error decoding encryption key", "err", base64Err.Error())
		panic("Error decoding encryption key")
	}

	encryptorManager, encryptorErr := utils.NewEncryptorManager(decodedSecret)

	if encryptorErr != nil {
		logger.Logger.Error("Error initializing encryptor manager", "err", encryptorErr.Error())
		panic("Error initializing encryptor manager")
	}

	appState := api.AppState{
		Db:               configs.InitDB(dbConfig),
		EncryptorManager: encryptorManager,
	}

	r := gin.Default()

	appState.SetupRoutes(r)
	appState.SetUpAuthRoutes(r)

	err := r.Run(fmt.Sprintf(":%s", settings.ServicePort))
	if err != nil {
		logger.Logger.Error("Error initializing service", "err", err.Error())
	}
}
