package config

import (
	"os"
	"github.com/joho/godotenv"
)

type Config struct {
	DBHost		string
	DBPort		string
	DBUser		string
	DBPassword	string
	DBName		string
	JWTSecret	string
	Port		string
}

func LoadConfig() *Config {
	godotenv.Load()

	return &Config{
		DBHost:		getEnv("DB_HOST", "127.0.0.1"),
		DBPort:		getEnv("DB_PORT", "3306"),
		DBUser:		getEnv("DB_USER", "root"),
		DBPassword:	getEnv("DB_PASSWORD", ""),
		DBName:		getEnv("DB_NAME", "db_go_ecommerce"),
		JWTSecret:	getEnv("JWT_SECRET", "bukankan_ini_my_secret_key_ku"),
		Port:		getEnv("PORT", "8080"),
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}

	return defaultValue
}