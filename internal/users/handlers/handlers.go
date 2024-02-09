package handlers

import (
	"bikesRentalAPI/internal/helpers"
	"bikesRentalAPI/internal/users/models"
	"bikesRentalAPI/internal/users/repository"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/go-chi/jwtauth/v5"
	"github.com/go-playground/validator/v10"
)

type Handler interface {
	LoginUser(tokenAuth *jwtauth.JWTAuth, w http.ResponseWriter, req *http.Request)
}

type handler struct {
	UserRepo  repository.UserRepository
	validator *validator.Validate
}

// New returns a new user handler
func New(userRepo repository.UserRepository) Handler {
	validator := validator.New(validator.WithRequiredStructEnabled())
	handler := &handler{
		UserRepo:  userRepo,
		validator: validator,
	}
	return handler
}

// RegisterUser ...
func RegisterUser(w http.ResponseWriter, req *http.Request) {
	// TODO Implement user registration logic

}

// LoginUser receives a tokenAuth and a request and returns a response
func (h *handler) LoginUser(tokenAuth *jwtauth.JWTAuth, w http.ResponseWriter, req *http.Request) {

	// Parse form data from request url-data encoded body
	err := req.ParseForm()
	if err != nil {
		http.Error(w, "Error parsing form", http.StatusBadRequest)
		return
	}

	credentials := models.LoginUserRequest{
		Email:    req.FormValue("email"),
		Password: req.FormValue("password"),
	}
	err = h.validator.Struct(credentials)
	if err != nil {
		errors := err.(validator.ValidationErrors)
		http.Error(w, fmt.Sprintf("Validation errors: %s", errors), http.StatusBadRequest)
		return
	}

	loggedUser, err := h.UserRepo.GetUserByEmail(strings.ToLower(credentials.Email))
	if err != nil {
		log.Printf("Error getting user by email: %v", err)
		http.Error(w, "Invalid username or password.", http.StatusUnauthorized)
		return
	}

	if !loggedUser.CheckPassword(credentials.Password) {
		http.Error(w, "Invalid username or password.", http.StatusUnauthorized)
		return
	}

	claimsMap := map[string]interface{}{
		"sub":       strconv.FormatInt(loggedUser.GetID(), 10),
		"exp":       time.Now().Add(time.Hour * 24 * 30).Unix(),
		"email":     loggedUser.GetEmail(),
		"firstName": loggedUser.GetFirstName(),
		"lastName":  loggedUser.GetLastName(),
	}

	_, tokenString, err := tokenAuth.Encode(claimsMap)
	if err != nil {
		log.Printf("Error encoding token: %v", err)
		http.Error(w, "Error encoding token", http.StatusInternalServerError)
		return
	}
	loginResponse := models.LoginUserResponse{
		Token: tokenString,
	}

	helpers.WriteJSON(w, http.StatusOK, loginResponse)

}

// LoginUser ...
func GetUserProfile(w http.ResponseWriter, r *http.Request) {
	// TODO Implement logic to retrieve user profile
	_, claims, _ := jwtauth.FromContext(r.Context())
	helpers.WriteJSON(w, http.StatusOK, claims)
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
