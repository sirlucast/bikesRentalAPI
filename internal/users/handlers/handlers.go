package handlers

import (
	"bikesRentalAPI/internal/helpers"
	"bikesRentalAPI/internal/users/models"
	"bikesRentalAPI/internal/users/repository"
	"encoding/json"
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
	RegisterUser(w http.ResponseWriter, req *http.Request)
	LoginUser(tokenAuth *jwtauth.JWTAuth, w http.ResponseWriter, req *http.Request)
	GetUserProfile(w http.ResponseWriter, req *http.Request)
	UpdateUserProfile(w http.ResponseWriter, req *http.Request)
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
func (h *handler) RegisterUser(w http.ResponseWriter, req *http.Request) {
	// Parse body from request
	body, err := helpers.ParseBody(req.Body)
	if err != nil {
		http.Error(w, "Error parsing body request", http.StatusBadRequest)
		return
	}
	var newUser models.CreateUserRequest
	err = json.Unmarshal(body, &newUser)
	if err != nil {
		http.Error(w, "Error unmarshalling body request", http.StatusBadRequest)
		return
	}
	// Validate user input
	err = h.validator.Struct(newUser)
	if err != nil {
		errors := err.(validator.ValidationErrors)
		http.Error(w, fmt.Sprintf("Validation errors: %s", errors), http.StatusBadRequest)
		return
	}
	id, err := h.UserRepo.CreateUser(newUser)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error creating user: %v", err), http.StatusInternalServerError)
		return
	}
	helpers.WriteJSON(w, http.StatusCreated, id)

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

// GetUserProfile ...
func (h *handler) GetUserProfile(w http.ResponseWriter, req *http.Request) {
	_, claims, err := jwtauth.FromContext(req.Context())
	if err != nil {
		http.Error(w, "Error getting user claims", http.StatusInternalServerError)
		return
	}
	userId, err := strconv.Atoi(claims["sub"].(string))
	if err != nil {
		fmt.Println("value is not a string")
	}
	user, err := h.UserRepo.GetUserByID(int64(userId))
	if err != nil {
		log.Printf("Error getting user: %v", err)
		http.Error(w, "Error getting user", http.StatusInternalServerError)
		return
	}
	helpers.WriteJSON(w, http.StatusOK, user)
}

// UpdateUserProfile ...
func (h *handler) UpdateUserProfile(w http.ResponseWriter, req *http.Request) {
	_, claims, err := jwtauth.FromContext(req.Context())
	if err != nil {
		http.Error(w, "Error getting user claims", http.StatusInternalServerError)
		return
	}
	body, err := helpers.ParseBody(req.Body)
	if err != nil {
		http.Error(w, "Error parsing body request", http.StatusBadRequest)
		return
	}
	var updateUserReq models.UpdateUserRequest
	err = json.Unmarshal(body, &updateUserReq)
	if err != nil {
		http.Error(w, "Error unmarshalling body request", http.StatusBadRequest)
		return
	}
	// Validate user input
	err = h.validator.Struct(updateUserReq)
	if err != nil {
		errors := err.(validator.ValidationErrors)
		http.Error(w, fmt.Sprintf("Validation errors: %s", errors), http.StatusBadRequest)
		return
	}

	// Get user from database
	userId, err := strconv.Atoi(claims["sub"].(string))
	if err != nil {
		log.Printf("Error getting user id from claims: %v", err)
		http.Error(w, "Not valid sub in claims", http.StatusUnauthorized)
	}

	user, err := h.UserRepo.GetUserByID(int64(userId))
	if err != nil {
		log.Printf("Error getting user: %v", err)
		http.Error(w, "Error getting user", http.StatusInternalServerError)
		return
	}

	// Compare attributes to identify modifications
	fieldsToUpdate, err := getFieldsToUpdate(&updateUserReq, &user)
	if err != nil {
		http.Error(w, "No fields to update", http.StatusBadRequest)
		return
	}

	// Update user
	result, err := h.UserRepo.UpdateUser(userId, fieldsToUpdate)
	if err != nil {
		http.Error(w, "Error updating user", http.StatusInternalServerError)
		return
	}
	updateUserResp := models.UpdateUserResponse{
		UserID:  result,
		Message: "User updated successfully",
	}
	helpers.WriteJSON(w, http.StatusOK, updateUserResp)
}

// getFieldsToUpdate ...
func getFieldsToUpdate(updateUserReq *models.UpdateUserRequest, user *models.User) (map[string]interface{}, error) {
	if updateUserReq == nil || user == nil {
		return nil, fmt.Errorf("no fields to update")
	}
	fieldsToUpdate := make(map[string]interface{})
	if updateUserReq.Email != nil && *updateUserReq.Email != user.GetEmail() {
		fieldsToUpdate["email"] = *updateUserReq.Email
	}
	if updateUserReq.FirstName != nil && *updateUserReq.FirstName != user.GetFirstName() {
		fieldsToUpdate["first_name"] = *updateUserReq.FirstName
	}
	if updateUserReq.LastName != nil && *updateUserReq.LastName != user.GetLastName() {
		fieldsToUpdate["last_name"] = *updateUserReq.LastName
	}

	if len(fieldsToUpdate) == 0 {
		return nil, fmt.Errorf("no fields to update")
	}
	return fieldsToUpdate, nil
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
