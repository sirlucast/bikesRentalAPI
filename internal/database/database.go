package database

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	_ "github.com/mattn/go-sqlite3"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/sqlite3"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

var dbUrl = os.Getenv("DB_URL")

type Database interface {
	Start() error
	Migrate() error
	QueryRow(string, ...interface{}) *sql.Row
	Exec(string, ...interface{}) (sql.Result, error)
	Prepare(string) (*sql.Stmt, error)
	Close() error
	Health() error
}

type database struct {
	db *sql.DB
}

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
	if err != nil {
		return fmt.Errorf("failed to create driver: %v", err)
	}
	// Create a new migration
	m, err := migrate.NewWithDatabaseInstance("file://internal/database/migrations", dbUrl, driver)

	if err != nil {
		return fmt.Errorf("failed to create migration: %v", err)
	}
	// Run the migrations
	err = m.Up()
	if err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("failed to run migrations: %v", err)
	}
	log.Println("Database migration completed successfully")
	return nil
}

func (d *database) Exec(query string, args ...interface{}) (sql.Result, error) {
	return d.db.Exec(query, args...)
}

func (d *database) Query(query string, args ...interface{}) (*sql.Rows, error) {
	return d.db.Query(query, args...)
}

func (d *database) QueryRow(query string, args ...interface{}) *sql.Row {
	return d.db.QueryRow(query, args...)
}

func (d *database) Prepare(query string) (*sql.Stmt, error) {
	return d.db.Prepare(query)
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
