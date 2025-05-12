package config

import (
	"log"
	"os"
	"strconv"
	"time"
)

type Config struct {
	AuthServiceAddr         string
	ServerPortAuth          string
	ServerPortChat          string
	NatsUrl                 string
	DBHost                  string
	DBHostUsers             string
	DBPort                  string
	DBPortUsers             string
	DBUser                  string
	DBUserUsers             string
	DBPassword              string
	DBPasswordUsers         string
	DBName                  string
	DBNameUsers             string
	SmtpUser                string
	SmtpPass                string
	SmtpHost                string
	SmtpPort                string
	JWTSecret               string
	AccessTokenDuration     time.Duration
	RefreshTokenDuration    time.Duration
	SagaPort                string
	NotificationServiceAddr string
	NotificationPort        string
}

func LoadConfig() *Config {
	return &Config{
		AuthServiceAddr: getEnv("AUTH_SERVICE_ADDR", "auth-service:50051"),
		ServerPortAuth:  getEnv("SERVER_PORT_AUTH", "50051"),
		ServerPortChat:  getEnv("SERVER_PORT_CHAT", "50052"),
		NatsUrl:         getEnv("NATS_URL", "nats://nats:4222"),

		DBHost:     getEnv("DB_HOST", "localhost"),
		DBPort:     getEnv("DB_PORT", "5432"),
		DBUser:     getEnv("DB_USER", "auth_user"),
		DBPassword: getEnv("DB_PASSWORD", "auth_pass"),
		DBName:     getEnv("DB_NAME", "auth_db"),

		DBHostUsers:     getEnv("DB_HOST_USERS", "localhost"),
		DBPortUsers:     getEnv("DB_PORT_USERS", "5432"),
		DBUserUsers:     getEnv("DB_USER_USERS", "user"),
		DBPasswordUsers: getEnv("DB_PASSWORD_USERS", "user_pass"),
		DBNameUsers:     getEnv("DB_NAME_USERS", "users_db"),

		SmtpHost: getEnv("SMTP_HOST", "smtp.gmail.com"),
		SmtpPort: getEnv("SMTP_PORT", "587"),
		SmtpUser: getEnv("SMTP_USER", "user@gmail.com"),
		SmtpPass: getEnv("SMTP_PASS", "pass-user"),

		JWTSecret:            getEnv("JWT_SECRET", "default_key"),
		AccessTokenDuration:  getEnvAsDuration("ACCESS_TOKEN_DURATION", time.Minute*15),
		RefreshTokenDuration: getEnvAsDuration("REFRESH_TOKEN_DURATION", time.Hour*24),

		SagaPort:                getEnv("SAGA_PORT", "50053"),
		NotificationPort:        getEnv("NOTIFICATION_PORT", "50054"),
		NotificationServiceAddr: getEnv("NOTIFICATION_SERVICE_ADDR", "notification-service:50054"),
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
