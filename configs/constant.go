package configs

import "os"

type Settings struct {
	DbPort      string
	DbHost      string
	DBName      string
	DbUser      string
	DbPassword  string
	ServicePort string
	ServiceName string
}

func getEnvOrDefault(key string, defaultValue string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return defaultValue
}

func GetSettings() *Settings {
	return &Settings{
		DbPort:      getEnvOrDefault("DB_PORT", "5432"),
		DbHost:      getEnvOrDefault("DB_HOST", "localhost"),
		DBName:      getEnvOrDefault("DB_NAME", "wordlee"),
		DbUser:      getEnvOrDefault("DB_USER", "postgres"),
		DbPassword:  getEnvOrDefault("DB_PASSWORD", "password"),
		ServicePort: getEnvOrDefault("SERVICE_PORT", "8080"),
		ServiceName: getEnvOrDefault("SERVICE_NAME", "wordlee-app-backend"),
	}
}
