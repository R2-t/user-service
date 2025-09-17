package configs

import (
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"santiagotorres.me/user-service/logger"
	"santiagotorres.me/user-service/models"
)

type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
}

func InitDB(cfg DatabaseConfig) *gorm.DB {

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s", cfg.Host, cfg.User, cfg.Password, cfg.DBName, cfg.Port, "disable")

	db, err := gorm.Open(postgres.New(postgres.Config{
		DSN: dsn,
	}), &gorm.Config{})

	if err != nil {
		logger.Logger.Error("Error connection to database")
		panic(err)
	}

	if err := db.AutoMigrate(&models.User{}); err != nil {
		logger.Logger.Error("Failed to migrate models", "err", err.Error())
		panic(err)
	}

	sqlDb, err := db.DB()

	if err != nil {
		logger.Logger.Error("Failed to get sql db", "err", err.Error())
		panic(err)
	}

	sqlDb.SetMaxIdleConns(20)

	return db
}
