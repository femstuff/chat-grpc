package repository

import (
	"database/sql"
	"fmt"
	"log"
)

func NewDb() (*sql.DB, error) {
	dsn := "host=localhost user=postgres password=1111 dbname=authdb sslmode=disable"

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("error connecting to db: %w", err)
	}

	if err = db.Ping(); err != nil {
		return nil, fmt.Errorf("error pinging db: %w", err)
	}

	log.Println("Connected to DB successfully")

	query := `
	CREATE TABLE IF NOT EXISTS users (
		id SERIAL PRIMARY KEY,
		name VARCHAR(255) NOT NULL,
		email VARCHAR(255) NOT NULL UNIQUE,
		password_hash TEXT NOT NULL,
		role VARCHAR(5) NOT NULL CHECK (role IN ('admin', 'user')),
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);`

	_, err = db.Exec(query)
	if err != nil {
		return nil, fmt.Errorf("error creating table: %w", err)
	}

	log.Println("Table 'users' is ready")

	return db, nil
}
