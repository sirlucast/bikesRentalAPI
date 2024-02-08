package database

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
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
	username, password, err := adminCredentialsDecode()
	if err != nil {
		return fmt.Errorf("failed Seed admin: %v", err)
	}
	queryString := fmt.Sprintf("SELECT * FROM users WHERE email = ?")
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

func adminCredentialsDecode() (string, string, error) {
	adminCredentials := helpers.SafeGetEnv("ADMIN_CREDENTIALS")
	decodedAdminCred, err := helpers.Base64Decode(adminCredentials)
	if err {
		return "", "", fmt.Errorf("failed to decode admin credentials.")
	}
	credentials := strings.Split(decodedAdminCred, ":")
	if len(credentials) != 2 {
		return "", "", fmt.Errorf("decoded admin credentials are not following <user:passowrd> shape. Got: %v", credentials)
	}
	return credentials[0], credentials[1], nil
}
