package models

import (
	"bikesRentalAPI/internal/helpers"
	"fmt"
	"time"
)

// User model
type User struct {
	ID             int64     `json:"id,omitempty"`
	Email          string    `json:"email,omitempty"`
	HashedPassword string    `json:"password,omitempty"`
	FirstName      *string   `json:"first_name,omitempty"`
	LastName       *string   `json:"last_name,omitempty"`
	CreatedAt      time.Time `json:"created_at,omitempty"`
	UpdatedAt      time.Time `json:"updated_at,omitempty"`
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
	Email    string `json:"email" validate:"required,email,omitempty"`
	Password string `json:"password" validate:"required"`
}

type LoginUserResponse struct {
	Token string
}

type CreateUserRequest struct {
	Email     string `json:"email" validate:"required,email,omitempty"`
	Password  string `json:"password" validate:"required"`
	FirstName string `json:"first_name" validate:"required,max=50"`
	LastName  string `json:"last_name" validate:"required,max=50"`
}

type UpdateUserRequest struct {
	Email     *string `json:"email" validate:"omitempty,email"`
	FirstName *string `json:"first_name" validate:"omitempty,max=50"`
	LastName  *string `json:"last_name" validate:"omitempty,max=50"`
	//NewPassowrd string  `json:"new_password" validate:"omitempty"`
	//OldPassword string  `json:"old_password" validate:"omitempty"`
}

type UpdateUserResponse struct {
	UserID  int64
	Message string
}

type Claims struct {
	Sub       int64  `json:"sub"`
	Exp       int64  `json:"exp"`
	Email     string `json:"email"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}
