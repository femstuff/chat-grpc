package pkg

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
	"go.uber.org/zap"
)

func NewDb(log *zap.Logger) (*sql.DB, error) {
	dsn := "host=localhost user=postgres password=1111 dbname=authdb sslmode=disable"

	db, err := sql.Open("postgres", dsn)
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
