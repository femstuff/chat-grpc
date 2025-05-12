package pkg

import (
	"database/sql"
	"fmt"

	"chat-grpc/pkg/config"
	_ "github.com/lib/pq"
	"go.uber.org/zap"
)

func NewDbChat(log *zap.Logger) (*sql.DB, error) {
	cfg := config.LoadConfig()
	connStr := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPassword, cfg.DBName,
	)
	log.Info("Connecting to chat DB", zap.String("connStr", connStr))

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

func NewDbUsers(log *zap.Logger) (*sql.DB, error) {
	cfg := config.LoadConfig()
	connStr := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.DBHostUsers, cfg.DBPortUsers, cfg.DBUserUsers, cfg.DBPasswordUsers, cfg.DBNameUsers)
	log.Info("Connecting to user DB", zap.String("connStr", connStr))

	db, err := sql.Open("postgres", connStr)

	if err != nil {
		log.Error("Failed connecting to db users", zap.Error(err))
		return nil, fmt.Errorf("error connect to db users")
	}

	if err = db.Ping(); err != nil {
		log.Error("Failed pinging users db")
		return nil, fmt.Errorf("error pingigng users db")
	}

	log.Info("Connected to DB users successfully")
	return db, nil
}
