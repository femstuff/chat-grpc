package config

import (
	"log"
	"os"
	"strconv"
	"time"
)

type Config struct {
	AUTH_SERVICE_ADDR    string
	SERVER_PORT_AUTH     string
	SERVER_PORT_CHAT     string
	DBHost               string
	DBPort               string
	DBUser               string
	DBPassword           string
	DBName               string
	JWTSecret            string
	AccessTokenDuration  time.Duration
	RefreshTokenDuration time.Duration
}

func LoadConfig() *Config {
	return &Config{
		AUTH_SERVICE_ADDR:    getEnv("AUTH_SERVICE_ADDR", "auth-service:50051"),
		SERVER_PORT_AUTH:     getEnv("SERVER_PORT_AUTH", "50051"),
		SERVER_PORT_CHAT:     getEnv("SERVER_PORT_CHAT", "50052"),
		DBHost:               getEnv("DB_HOST", "localhost"),
		DBPort:               getEnv("DB_PORT", "5432"),
		DBUser:               getEnv("DB_USER", "auth_user"),
		DBPassword:           getEnv("DB_PASSWORD", "auth_pass"),
		DBName:               getEnv("DB_NAME", "auth_db"),
		JWTSecret:            getEnv("JWT_SECRET", "default_key"),
		AccessTokenDuration:  getEnvAsDuration("ACCESS_TOKEN_DURATION", time.Minute*15),
		RefreshTokenDuration: getEnvAsDuration("REFRESH_TOKEN_DURATION", time.Hour*24),
	}
}

func getEnv(key, defaultVal string) string {
	val := os.Getenv(key)
	if val == "" {
		return defaultVal
	}

	return val
}

func getEnvAsDuration(key string, defaultVal time.Duration) time.Duration {
	valStr := os.Getenv(key)
	if valStr == "" {
		return defaultVal
	}
	valInt, err := strconv.Atoi(valStr)
	if err != nil {
		log.Printf("Invalid duration for %s: %s, using default\n", key, valStr)
		return defaultVal
	}

	return time.Duration(valInt) * time.Second
}
