package handlers

import (
	"bikesRentalAPI/internal/helpers"
	"bikesRentalAPI/internal/users/models"
	"bikesRentalAPI/internal/users/repository"
	"encoding/json"
	"io"
	"log"
	"net/http"
)

type Handler struct {
	UserRepo repository.UserRepository
}

// New returns a new user handler
func New(userRepo repository.UserRepository) *Handler {
	return &Handler{UserRepo: userRepo}
}

// RegisterUser ...
func RegisterUser(w http.ResponseWriter, req *http.Request) {
	// TODO Implement user registration logic
}

// LoginUser ...
func (h *Handler) LoginUser(w http.ResponseWriter, r *http.Request) {
	// TODO complete this
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Error reading request body",
			http.StatusInternalServerError)
		return
	}
	log.Printf("Request body: %s\n", body)
	var creds models.LoginUserRequest
	err = json.Unmarshal(body, &creds)
	if err != nil {
		log.Printf("Error unmarshalling request body: %v", err)
		http.Error(w, "Error unmarshalling request body",
			http.StatusBadRequest)
		return
	}
	if creds.Email == "" || creds.Password == "" {
		http.Error(w, "Missing username or password.", http.StatusBadRequest)
		return
	}
	testUser, _ := h.UserRepo.GetUserByEmail(creds.Email)
	helpers.WriteJSON(w, http.StatusOK, testUser)

}

// LoginUser ...
func GetUserProfile(w http.ResponseWriter, r *http.Request) {
	// TODO Implement logic to retrieve user profile
}

// UpdateUserProfile ...
func UpdateUserProfile(w http.ResponseWriter, r *http.Request) {
	// TODO Implement logic to update user profile
}

// Admin access only

// ListUsers ...
func ListUsers(w http.ResponseWriter, r *http.Request) {
	// TODO Implement logic to list all users
}

// GetUserDetails ...
func GetUserDetails(w http.ResponseWriter, r *http.Request) {
	// TODO Implement logic to get user details
}

// UpdateUserDetails ...
func UpdateUserDetails(w http.ResponseWriter, r *http.Request) {
	// TODO Implement logic to update user details
}
