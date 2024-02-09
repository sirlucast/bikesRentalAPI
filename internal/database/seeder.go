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
	Seed(users.User) error
}

type seeder struct {
	Database Database
}

func NewSeeder(db Database) Seeder {
	return &seeder{Database: db}
}

// Seed populates the database with initial data
func (s *seeder) Seed(user users.User) error {
	username, password, err := userCredentialsDecode()
	username = strings.ToLower(username)
	if err != nil {
		return fmt.Errorf("failed Seed admin: %v", err)
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
				return fmt.Errorf("failed to insert admin user: %v", err)
			}
			log.Printf("admin user succesfully seeded")
			return nil
		}
	}
	log.Printf("admin user already exists")
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
