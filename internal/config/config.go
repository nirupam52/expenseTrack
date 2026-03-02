package config

import "os"

type DBConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	Database string
	DBPath   string
}

type AppConfig struct {
	Port      string
	DbConfig  DBConfig
	SecretKey string
}

func LoadConfig() AppConfig {
	return AppConfig{
		Port:      getEnvVar("APP_PORT", "41605"),
		DbConfig:  LoadDBConfig(),
		SecretKey: os.Getenv("APP_SECRET_KEY"),
	}
}

func LoadDBConfig() DBConfig {
	return DBConfig{
		Host:     os.Getenv("DB_HOST"),
		Port:     5432,
		User:     os.Getenv("DB_USER"),
		Password: os.Getenv("DB_PASSWORD"),
		Database: os.Getenv("DB_NAME"),
		DBPath: getEnvVar("DB_PATH", "expense.db"),
	}
}

func getEnvVar(key string, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}
	return value
}
