package handlers

import (
	"bikesRentalAPI/internal/helpers"
	"bikesRentalAPI/internal/rentals/models"
	"bikesRentalAPI/internal/rentals/repository"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/go-chi/jwtauth/v5"
	"github.com/go-playground/validator/v10"
)

// Handler is the interface for rental handlers
type Handler interface {
	GetRentalHistoryByUserID(w http.ResponseWriter, req *http.Request) // Get rental history by user ID
	GetRentalList(w http.ResponseWriter, req *http.Request)            // Get rental list
	GetRentalDetails(w http.ResponseWriter, req *http.Request)         // Get rental details
	UpdateRentalDetails(w http.ResponseWriter, req *http.Request)      // Update rental details
	StartBikeRental(w http.ResponseWriter, req *http.Request)          // Start bike rental
	EndBikeRental(w http.ResponseWriter, req *http.Request)            // End bike rental
}

type handler struct {
	RentalRepo repository.RentalRepository
	validator  *validator.Validate
}

// New returns a new rental handler
func New(RentalRepository repository.RentalRepository) Handler {
	validator := validator.New(validator.WithRequiredStructEnabled())
	handler := &handler{
		RentalRepo: RentalRepository,
		validator:  validator,
	}
	return handler
}

// GetRentalList ...
func (h *handler) GetRentalList(w http.ResponseWriter, req *http.Request) {
	// TODO Implement logic to list all rentals
}

// GetRentalDetails ...
func (h *handler) GetRentalDetails(w http.ResponseWriter, req *http.Request) {
	// TODO Implement logic to get rental details
}

// UpdateRentalDetails ...
func (h *handler) UpdateRentalDetails(w http.ResponseWriter, req *http.Request) {
	// TODO Implement logic to update rental details
}

// StartBikeRental ...
func (h *handler) StartBikeRental(w http.ResponseWriter, req *http.Request) {
	_, claims, err := jwtauth.FromContext(req.Context())
	if err != nil || claims == nil {
		http.Error(w, "Error getting user claims", http.StatusBadRequest)
		return
	}
	userId, err := strconv.ParseInt(claims["sub"].(string), 10, 64)
	if err != nil {
		log.Printf("Error getting user id from jwt.claims: %v", err)
		http.Error(w, "Error getting user id", http.StatusBadRequest)
		return
	}
	isUserRenting := h.RentalRepo.IsUserRentingBike(userId)
	if isUserRenting {
		http.Error(w, "User is already renting a bike", http.StatusBadRequest)
		return
	}

	body, err := helpers.ParseBody(req.Body)
	if err != nil {
		http.Error(w, "Error parsing body request", http.StatusBadRequest)
		return
	}

	var startBikeRentalReq *models.StartBikeRentalRequest
	if err := json.Unmarshal(body, &startBikeRentalReq); err != nil {
		http.Error(w, "Error decoding body request", http.StatusBadRequest)
		return
	}
	if err := h.validator.Struct(startBikeRentalReq); err != nil {
		errors := err.(validator.ValidationErrors)
		http.Error(w, fmt.Sprintf("Validation errors: %s", errors), http.StatusBadRequest)
		return
	}

	isBikeAvailable := h.RentalRepo.IsBikeAvailable(startBikeRentalReq.BikeID)
	if !isBikeAvailable {
		http.Error(w, "Bike is not available for rent", http.StatusBadRequest)
		return
	}

	rental, err := h.RentalRepo.StartRental(userId, startBikeRentalReq)
	if err != nil {
		log.Printf("Error starting bike rental: %v", err)
		http.Error(w, "Error starting bike rental", http.StatusBadRequest)
		return
	}

	helpers.WriteJSON(w, http.StatusOK, rental)
}

// EndBikeRental ...
func (h *handler) EndBikeRental(w http.ResponseWriter, req *http.Request) {
	_, claims, err := jwtauth.FromContext(req.Context())
	if err != nil || claims == nil {
		http.Error(w, "Error getting user claims", http.StatusBadRequest)
		return
	}
	userId, err := strconv.ParseInt(claims["sub"].(string), 10, 64)
	if err != nil {
		log.Printf("Error getting user id from jwt.claims: %v", err)
		http.Error(w, "Error getting user id", http.StatusBadRequest)
		return
	}

	isUserRenting := h.RentalRepo.IsUserRentingBike(userId)
	if !isUserRenting {
		http.Error(w, "User is not currntly renting a bike", http.StatusBadRequest)
		return
	}

	body, err := helpers.ParseBody(req.Body)
	if err != nil {
		http.Error(w, "Error parsing body request", http.StatusBadRequest)
		return
	}

	var startBikeRentalReq *models.StopBikeRentalRequest
	if err := json.Unmarshal(body, &startBikeRentalReq); err != nil {
		http.Error(w, "Error decoding body request", http.StatusBadRequest)
		return
	}
	if err := h.validator.Struct(startBikeRentalReq); err != nil {
		errors := err.(validator.ValidationErrors)
		http.Error(w, fmt.Sprintf("Validation errors: %s", errors), http.StatusBadRequest)
		return
	}

	ongoingRental, err := h.RentalRepo.GetOngoingRental(userId)
	if err != nil {
		log.Printf("Error getting ongoing rental: %v", err)
		http.Error(w, "Error getting ongoing rental", http.StatusBadRequest)
		return
	}
	if ongoingRental.ID != startBikeRentalReq.RentalID {
		http.Error(w, "Rental in request does not match with current bike rental by user", http.StatusBadRequest)
		return
	}

	rental, err := h.RentalRepo.EndRental(userId, startBikeRentalReq)
	if err != nil {
		log.Printf("Error ending bike rental: %v", err)
		http.Error(w, "Error ending bike rental", http.StatusBadRequest)
		return
	}
	helpers.WriteJSON(w, http.StatusOK, rental)

}

// GetRentalHistoryByUserID ...
func (h *handler) GetRentalHistoryByUserID(w http.ResponseWriter, req *http.Request) {
	// TODO Implement logic to get rental history
	// 1. Get pagination details from the request context
	//pageID := r.Context().Value(middlewares.PageIDKey)
	// 2. Get rental history from the database
}
