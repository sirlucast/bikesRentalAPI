package database

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"time"

	_ "github.com/joho/godotenv/autoload"
	_ "github.com/mattn/go-sqlite3"
)

type Database interface {
	Health() error
}

type database struct {
	db *sql.DB
}

var (
	dbUrl = os.Getenv("DB_URL")
)

// New initializes a new database service
func New() (Database, error) {
	db, err := sql.Open("sqlite3", dbUrl)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return &database{db: db}, nil
}

// Health checks the database connection and returns an error if it's down
func (s *database) Health() error {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	if err := s.db.PingContext(ctx); err != nil {
		return fmt.Errorf("database down: %w", err)
	}

	return nil
}
