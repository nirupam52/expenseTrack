package config

import "os"

type AppConfig struct {
	Port      string
	DBPath    string
	SecretKey string
}

func LoadConfig() AppConfig {
	return AppConfig{
		Port:      getEnvVar("APP_PORT", "41605"),
		DBPath:    getEnvVar("DB_PATH", "expense.db"),
		SecretKey: os.Getenv("APP_SECRET_KEY"),
	}
}

func getEnvVar(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
