package models

import (
	"bikesRentalAPI/internal/helpers"
	"fmt"
	"time"
)

// User model
type User struct {
	ID             int64
	Email          string
	HashedPassword string
	FirstName      *string
	LastName       *string
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

func (u *User) FullName() string {
	return *u.FirstName + " " + *u.LastName
}

func (u *User) SetPassword(password string) error {
	hash, err := helpers.GetHashPassword(password)
	if err != nil {
		return fmt.Errorf("failed to hash password: %v", err)
	}
	u.HashedPassword = hash
	return nil
}

func (u *User) CheckPassword(password string) bool {
	return helpers.CheckPassword(u.HashedPassword, password)
}

func (u *User) GetID() int64 {
	return u.ID
}

func (u *User) GetFirstName() string {
	if u.FirstName == nil {
		return ""
	}
	return *u.FirstName
}

func (u *User) GetLastName() string {
	if u.LastName == nil {
		return ""
	}
	return *u.LastName
}

func (u *User) GetEmail() string {
	return u.Email
}

type LoginUserRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginUserResponse struct {
	Token string
}

type CreateUserRequest struct {
	Email     string `json:"email"`
	Password  string `json:"password"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

type CreateUserResponse struct {
	ID        int64  `json:"id"`
	Email     string `json:"email"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

type UpdateUserRequest struct {
	ID        int64  `json:"id"`
	Email     string `json:"email"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Passowrd  string `json:"password"`
}

type Claims struct {
	Sub       int64  `json:"sub"`
	Exp       int64  `json:"exp"`
	Email     string `json:"email"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}
