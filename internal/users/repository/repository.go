package repository

import (
	"bikesRentalAPI/internal/database"
	"bikesRentalAPI/internal/helpers"
	"bikesRentalAPI/internal/users/models"
	"fmt"
	"strings"
)

const (
	// pageSize is the number of items to return in a page, 10 as default.
	pageSize = 10
)

type UserRepository interface {
	CreateUser(models.CreateUserRequest) (int64, error)
	GetUserByEmailForAuth(string) (*models.User, error)
	GetUserByID(int64) (*models.User, error)
	UpdateUser(userID int64, fieldsToUpdate map[string]interface{}) (int64, error)
	ListAllUsers(int64) (*models.UserList, error)
	IsEmailUnique(string) (bool, error)
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

// GetUserByEmailForAuth retrieves a user from the database by email
func (r *userRepository) GetUserByEmailForAuth(email string) (*models.User, error) {
	var user models.User
	query := "SELECT id, email, hashed_password, first_name, last_name FROM users WHERE email = ?"
	err := r.db.QueryRow(query, email).Scan(&user.ID, &user.Email, &user.HashedPassword, &user.FirstName, &user.LastName)
	if err != nil {
		return &user, err
	}
	return &user, nil
}

func (r *userRepository) IsEmailUnique(email string) (bool, error) {
	var count int
	query := "SELECT COUNT(*) FROM users WHERE email = ?"
	err := r.db.QueryRow(query, email).Scan(&count)
	if err != nil {
		return false, err
	}
	return count == 0, nil
}

// GetUserByID retrieves a user from the database by id
func (r *userRepository) GetUserByID(id int64) (*models.User, error) {
	var user models.User
	quwery := "SELECT id, email, first_name, last_name, created_at, updated_at FROM users WHERE id = ?"
	err := r.db.QueryRow(quwery, id).Scan(&user.ID, &user.Email, &user.FirstName, &user.LastName, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		return &user, err
	}
	return &user, nil
}

// UpdateUser updates a user in the database by id. Returns the id of the updated user. If no fields are updated, returns 0 and an error
func (r *userRepository) UpdateUser(userID int64, fieldsToUpdate map[string]interface{}) (int64, error) {
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
	if err != nil {
		return 0, fmt.Errorf("failed to get last insert id: %v", err)
	}
	return id, nil
}

// ListAllUsers retrieves all bikes from the database
func (u *userRepository) ListAllUsers(PageID int64) (*models.UserList, error) {
	query := "SELECT id, email, first_name, last_name, created_at, updated_at FROM bikes WHERE id > ? ORDER BY id LIMIT ?"
	rows, err := u.db.Query(query, PageID, pageSize)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	users := &models.UserList{}
	userList := make([]*models.User, 0)
	for rows.Next() {
		var user models.User
		if err := rows.Scan(&user.ID, &user.Email, &user.FirstName, &user.LastName, &user.CreatedAt, &user.UpdatedAt); err != nil {
			return nil, err
		}
		userList = append(userList, &user)
	}
	if len(userList) == pageSize {
		users.NextPageID = userList[len(userList)-1].ID
	}
	users.Items = userList
	return users, nil
}
