package database

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"os"
	"strings"

	"bikesRentalAPI/internal/helpers"
	users "bikesRentalAPI/internal/users/models"
)

type Seeder interface {
	SeedUser() error
	SeedBikes() error
}

type seeder struct {
	Database Database
}

func NewSeeder(db Database) Seeder {
	return &seeder{Database: db}
}

// Seed populates the database with initial data
func (s *seeder) SeedUser() error {
	user := users.User{}
	username, password, err := userCredentialsDecode()
	username = strings.ToLower(username)
	if err != nil {
		return fmt.Errorf("failed Seed user: %v", err)
	}
	queryString := "SELECT * FROM users WHERE email = ?"
	row := s.Database.QueryRow(queryString, username)
	if err := row.Scan(&user.Email); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			hashedPassword, err := helpers.GetHashPassword(password)
			if err != nil {
				return fmt.Errorf("failed to hash password: %v", err)
			}
			_, err = s.Database.Exec("INSERT INTO users (email, hashed_password) VALUES (?, ?)", username, hashedPassword)
			if err != nil {
				return fmt.Errorf("failed to insert user: %v", err)
			}
			log.Printf("user succesfully seeded")
			return nil
		}
	}
	log.Printf("user already exists")
	return nil
}

func (s *seeder) SeedBikes() error {
	PricePerMinute := 0.07
	latitude := 51.5098387087398
	longitude := -0.1626587921593317
	queryString := "INSERT INTO bikes (is_available, price_per_minute, latitude, longitude) VALUES (?, ?, ?, ?)"
	// Insert 10 bikes
	for i := 0; i < 10; i++ {
		_, err := s.Database.Exec(queryString, true, PricePerMinute, latitude, longitude)
		if err != nil {
			return fmt.Errorf("failed to insert bike: %v", err)
		}
	}
	log.Printf("bikes succesfully seeded")
	return nil
}

func userCredentialsDecode() (string, string, error) {
	envCreds := os.Getenv("USER_CREDENTIALS")
	credentials := strings.Split(envCreds, ":")
	if len(credentials) != 2 {
		return "", "", fmt.Errorf("decoded admin credentials are not following <user:passowrd> shape. Got: %v", credentials)
	}
	return credentials[0], credentials[1], nil
}
