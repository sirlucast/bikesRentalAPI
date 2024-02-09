package repository

import (
	"bikesRentalAPI/internal/database"
	"bikesRentalAPI/internal/helpers"
	"bikesRentalAPI/internal/users/models"
	"fmt"
	"log"
	"strings"
)

type UserRepository interface {
	CreateUser(models.CreateUserRequest) (int64, error)
	GetUserByEmail(string) (models.User, error)
	GetUserByID(int64) (models.User, error)
	UpdateUser(userID int, fieldsToUpdate map[string]interface{}) (int64, error)
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
	query := "SELECT id, email, first_name, last_name, created_at, updated_at FROM users WHERE email = ?"
	err := r.db.QueryRow(query, email).Scan(&user.ID, &user.Email, &user.HashedPassword, &user.FirstName, &user.LastName, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		return user, err
	}
	return user, nil
}

// GetUserByID retrieves a user from the database by id
func (r *userRepository) GetUserByID(id int64) (models.User, error) {
	user := models.User{}
	quwery := "SELECT id, email, first_name, last_name, created_at, updated_at FROM users WHERE id = ?"
	err := r.db.QueryRow(quwery, id).Scan(&user.ID, &user.Email, &user.FirstName, &user.LastName, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		return user, err
	}
	return user, nil
}

// UpdateUser updates a user in the database by id. Returns the id of the updated user. If no fields are updated, returns 0 and an error
func (r *userRepository) UpdateUser(userID int, fieldsToUpdate map[string]interface{}) (int64, error) {
	var setFields []string
	var args []interface{}

	for field, value := range fieldsToUpdate {
		setFields = append(setFields, fmt.Sprintf("%s = ?", field))
		args = append(args, value)
	}
	args = append(args, userID)

	query := fmt.Sprintf("UPDATE users SET %s WHERE id = ?", strings.Join(setFields, ", "))
	stmt, err := r.db.Prepare(query)
	if err != nil {
		return 0, fmt.Errorf("failed to prepare update statement: %v", err)
	}
	defer stmt.Close()

	result, err := stmt.Exec(args...)
	if err != nil {
		return 0, fmt.Errorf("failed to execute update statement: %v", err)
	}

	id, err := result.LastInsertId()
	log.Printf("id: %v", id)
	if err != nil {
		return 0, fmt.Errorf("failed to get last insert id: %v", err)
	}
	return id, nil
}
