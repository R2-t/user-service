package main

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"santiagotorres.me/user-service/api"
	"santiagotorres.me/user-service/configs"
	"santiagotorres.me/user-service/logger"
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

	appState := api.AppState{
		Db: configs.InitDB(dbConfig),
	}

	r := gin.Default()

	appState.SetupRoutes(r)
	appState.SetUpAuthRoutes(r)

	err := r.Run(fmt.Sprintf(":%s", settings.ServicePort))
	if err != nil {
		logger.Logger.Error("Error initializing service", "err", err.Error())
	}
}
