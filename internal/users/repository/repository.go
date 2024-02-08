package repository

import (
	"bikesRentalAPI/internal/database"
	"bikesRentalAPI/internal/helpers"
	"bikesRentalAPI/internal/users/models"
	"fmt"
)

type UserRepository interface {
	CreateUser(models.CreateUserRequest) (int64, error)
	GetUserByEmail(string) (models.User, error)
	GetUserByID(int64) (models.User, error)
	UpdateUser(models.User) error
}

type userRepository struct {
	db database.Database
}

// New initializes a new empty user repository
func New(db database.Database) UserRepository {
	return &userRepository{db}
}

// CreateUser inserts a new user into the database
func (r *userRepository) CreateUser(user models.CreateUserRequest) (int64, error) {
	hashedPsw, err := helpers.GetHashPassword(user.Password)
	if err != nil {
		return 0, fmt.Errorf("failed to hash password: %v", err)
	}
	result, err := r.db.Exec("INSERT INTO users (email, hashed_password, first_name, last_name) VALUES (?, ?, ?, ?)", user.Email, hashedPsw, user.FirstName, user.LastName)
	if err != nil {
		return 0, fmt.Errorf("failed to insert user: %v", err)
	}
	id, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("failed to get last insert id: %v", err)
	}
	return id, nil
}

// GetUserByEmail retrieves a user from the database by email
func (r *userRepository) GetUserByEmail(email string) (models.User, error) {
	user := models.User{}
	err := r.db.QueryRow("SELECT * FROM users WHERE email = ?", email).Scan(&user.ID, &user.Email, &user.HashedPassword, &user.FirstName, &user.LastName, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		return user, err
	}
	return user, nil
}

// GetUserByID retrieves a user from the database by id
func (r *userRepository) GetUserByID(id int64) (models.User, error) {
	user := models.User{}
	err := r.db.QueryRow("SELECT * FROM users WHERE id = ?", id).Scan(&user.ID, &user.Email, &user.HashedPassword, &user.FirstName, &user.LastName, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		return user, err
	}
	return user, nil
}

func (r *userRepository) UpdateUser(user models.User) error {
	_, err := r.db.Exec("UPDATE users SET email = ?, first_name = ?, last_name = ? WHERE id = ?", user.Email, user.FirstName, user.LastName, user.ID)
	if err != nil {
		return fmt.Errorf("failed to update user: %v", err)
	}
	return nil
}
