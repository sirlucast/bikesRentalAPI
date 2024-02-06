package database

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"time"

	_ "github.com/joho/godotenv/autoload"
	_ "github.com/mattn/go-sqlite3"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/sqlite3"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

type Database interface {
	Start() error
	Migrate() error
	Close() error
	Health() error
}

type database struct {
	db *sql.DB
}

var (
	dbUrl = os.Getenv("DB_URL")
)

// New initializes a new empty database service
func New() Database {
	return &database{db: nil}
}

// Start initializes the database connection
func (d *database) Start() error {
	db, err := sql.Open("sqlite3", dbUrl)
	if err != nil {
		return fmt.Errorf("failed to open database: %w", err)
	}

	if err := db.Ping(); err != nil {
		return fmt.Errorf("failed to ping database: %w", err)
	}
	d.db = db
	return nil
}

func (d *database) Migrate() error {
	driver, err := sqlite3.WithInstance(d.db, &sqlite3.Config{})
	m, err := migrate.NewWithDatabaseInstance("file://internal/database/migrations", dbUrl, driver)

	if err != nil {
		return fmt.Errorf("failed to create migration: %v", err)
	}
	// Run the migrations
	err = m.Up()
	if err != nil {
		return fmt.Errorf("failed to run migrations: %v", err)
	}
	return nil
}

// Health checks the database connection and returns an error if it's down
func (d *database) Health() error {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	if err := d.db.PingContext(ctx); err != nil {
		return fmt.Errorf("database down: %w", err)
	}

	return nil
}

func (d *database) Close() error {
	if d.db != nil {
		return d.db.Close()
	}
	return nil
}
