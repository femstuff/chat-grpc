package pkg

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/lib/pq"
	"go.uber.org/zap"
)

func NewDb(log *zap.Logger) (*sql.DB, error) {
	connStr := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		getEnv("DB_HOST", "localhost"),
		getEnv("DB_PORT", "5432"),
		getEnv("DB_USER", "postgres"),
		getEnv("DB_PASSWORD", "password"),
		getEnv("DB_NAME", "postgres"),
	)
	log.Info("Connecting to DB")

	db, err := sql.Open("postgres", connStr)

	if err != nil {
		log.Error("Failed connecting to db", zap.Error(err))
		return nil, fmt.Errorf("error connecting to db: %w", err)
	}

	if err = db.Ping(); err != nil {
		log.Error("Failed pinging db", zap.Error(err))
		return nil, fmt.Errorf("error pinging db: %w", err)
	}

	log.Info("Connected to DB successfully")
	return db, nil
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}
