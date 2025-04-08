package pkg

import (
	"database/sql"
	"fmt"

	"chat-grpc/pkg/config"
	_ "github.com/lib/pq"
	"go.uber.org/zap"
)

func NewDb(log *zap.Logger) (*sql.DB, error) {
	cfg := config.LoadConfig()
	connStr := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPassword, cfg.DBName,
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
