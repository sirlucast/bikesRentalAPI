package models

import (
	"bikesRentalAPI/internal/helpers"
	"fmt"
	"time"
)

// User model represents a user of the application
type User struct {
	ID             int64     `json:"id,omitempty"`
	Email          string    `json:"email,omitempty"`
	HashedPassword string    `json:"password,omitempty"`
	FirstName      *string   `json:"first_name,omitempty"`
	LastName       *string   `json:"last_name,omitempty"`
	CreatedAt      time.Time `json:"created_at,omitempty"`
	UpdatedAt      time.Time `json:"updated_at,omitempty"`
}

// FullName returns the full name of the user
func (u *User) FullName() string {
	return u.GetFirstName() + " " + u.GetLastName()
}

// SetPassword hashes the password and sets it to the user
func (u *User) SetPassword(password string) error {
	hash, err := helpers.GetHashPassword(password)
	if err != nil {
		return fmt.Errorf("failed to hash password: %v", err)
	}
	u.HashedPassword = hash
	return nil
}

// CheckPassword checks if the password is correct
func (u *User) CheckPassword(password string) bool {
	return helpers.CheckPassword(u.HashedPassword, password)
}

// GetID returns the id of the user
func (u *User) GetID() int64 {
	return u.ID
}

// GetFirstName returns the first name of the user
func (u *User) GetFirstName() string {
	if u.FirstName == nil {
		return ""
	}
	return *u.FirstName
}

// GetLastName returns the last name of the user
func (u *User) GetLastName() string {
	if u.LastName == nil {
		return ""
	}
	return *u.LastName
}

// GetEmail returns the email of the user
func (u *User) GetEmail() string {
	return u.Email
}

// LoginUserRequest represents the request to login a user
type LoginUserRequest struct {
	Email    string `json:"email" validate:"required,email,omitempty"`
	Password string `json:"password" validate:"required"`
} // @name LoginUserRequest

// LoginUserResponse represents the response to login a user
type LoginUserResponse struct {
	Token string
} // @name LoginUserResponse

// CreateUserRequest represents the request to create a user
type CreateUserRequest struct {
	Email     string `json:"email" validate:"required,email,omitempty"`
	Password  string `json:"password" validate:"required"`
	FirstName string `json:"first_name" validate:"required,max=50"`
	LastName  string `json:"last_name" validate:"required,max=50"`
} // @name CreateUserRequest

type UpdateUserRequest struct {
	Email     *string `json:"email" validate:"omitempty,email"`
	FirstName *string `json:"first_name" validate:"omitempty,max=50"`
	LastName  *string `json:"last_name" validate:"omitempty,max=50"`
	//NewPassowrd string  `json:"new_password" validate:"omitempty"`
	//OldPassword string  `json:"old_password" validate:"omitempty"`
} // @name UpdateUserRequest

// UpdateUserResponse represents the response to update a user
type UpdateUserResponse struct {
	UserID  int64
	Message string
} // @name UpdateUserResponse

// Claims represents the claims of a JWT token
type Claims struct {
	Sub       int64  `json:"sub"`
	Exp       int64  `json:"exp"`
	Email     string `json:"email"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
} // @name Claims
