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

	query := `
	CREATE TABLE IF NOT EXISTS users (
		id SERIAL PRIMARY KEY,
		name VARCHAR(255) NOT NULL,
		email VARCHAR(255) NOT NULL UNIQUE,
		password_hash TEXT NOT NULL,
		role VARCHAR(5) NOT NULL CHECK (role IN ('admin', 'user')),
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);

	CREATE TABLE IF NOT EXISTS refresh_tokens (
		user_id BIGINT PRIMARY KEY,
		token TEXT NOT NULL
	);
	
	CREATE  TABLE  IF NOT EXISTS messages (
	  id SERIAL PRIMARY KEY,
	  chat_id INT NOT NULL,
	  user_id BIGINT NOT NULL,
	  text TEXT NOT NULL,
	  timestamp TIMESTAMP DEFAULT CURRENT_TIMESTAMP    
	);
	
	CREATE TABLE IF NOT EXISTS chats (
		id SERIAL PRIMARY KEY,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);
	CREATE TABLE IF NOT EXISTS chat_users (
    chat_id INT NOT NULL REFERENCES chats(id) ON DELETE CASCADE,
    user_id INT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    PRIMARY KEY (chat_id, user_id)
	);`

	_, err = db.Exec(query)
	if err != nil {
		log.Error("Failed creating table", zap.Error(err))
		return nil, fmt.Errorf("error creating table: %w", err)
	}

	log.Info("Tables is ready")

	return db, nil
}
