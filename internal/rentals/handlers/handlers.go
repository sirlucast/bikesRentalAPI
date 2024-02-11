package handlers

import (
	"bikesRentalAPI/internal/helpers"
	"bikesRentalAPI/internal/middlewares"
	"bikesRentalAPI/internal/rentals/models"
	"bikesRentalAPI/internal/rentals/repository"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
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
	pageID := req.Context().Value(middlewares.PageIDKey)

	rentals, err := h.RentalRepo.ListAllRentals(pageID.(int64))
	if err != nil {
		log.Printf("Error getting rental history: %v", err)
		http.Error(w, "Error getting rental history", http.StatusBadRequest)
		return
	}
	helpers.WriteJSON(w, http.StatusOK, rentals)
}

// GetRentalDetails ...
func (h *handler) GetRentalDetails(w http.ResponseWriter, req *http.Request) {
	rentalIDStr := chi.URLParam(req, "rental_id")
	rentalID, err := strconv.ParseInt(rentalIDStr, 10, 64)
	if err != nil {
		http.Error(w, fmt.Sprintf("couldn't read %s: %v", rentalIDStr, err), http.StatusBadRequest)
		return
	}

	rental, err := h.RentalRepo.GetRentalDetails(rentalID)
	if err != nil {
		log.Printf("Error getting rental details: %v", err)
		http.Error(w, "Error getting rental details", http.StatusBadRequest)
		return
	}
	helpers.WriteJSON(w, http.StatusOK, rental)
}

// UpdateRentalDetails ...
func (h *handler) UpdateRentalDetails(w http.ResponseWriter, req *http.Request) {
	rentalIDStr := chi.URLParam(req, "rental_id")
	rentalID, err := strconv.ParseInt(rentalIDStr, 10, 64)
	if err != nil {
		http.Error(w, fmt.Sprintf("couldn't read %s: %v", rentalIDStr, err), http.StatusBadRequest)
		return
	}

	// Parse body from request
	body, err := helpers.ParseBody(req.Body)
	if err != nil {
		http.Error(w, "Error parsing body request", http.StatusBadRequest)
		return
	}

	var updateRentalReq *models.UpdateRentalRequest
	err = json.Unmarshal(body, &updateRentalReq)
	if err != nil {
		http.Error(w, "Error unmarshalling body request", http.StatusBadRequest)
		return
	}

	// Validate user input
	err = h.validator.Struct(updateRentalReq)
	if err != nil {
		errors := err.(validator.ValidationErrors)
		http.Error(w, fmt.Sprintf("Validation errors: %s", errors), http.StatusBadRequest)
		return
	}

	rental, err := h.RentalRepo.GetRentalDetails(rentalID)
	if err != nil {
		log.Printf("UpdateRentalDetails: Error getting rental details: %v", err)
		http.Error(w, "Error getting rental details", http.StatusBadRequest)
		return
	}

	// Compare attributes to identify modifications
	fieldsToUpdate, err := getFieldsToUpdate(updateRentalReq, rental)
	if err != nil {
		http.Error(w, "No fields to update", http.StatusBadRequest)
		return
	}

	// A user can only rent one bike at a time
	if _, ok := fieldsToUpdate["user_id"]; ok {
		isUserRenting := h.RentalRepo.IsUserRentingBike(*updateRentalReq.UserID)
		if isUserRenting {
			http.Error(w, "User is already renting a bike", http.StatusBadRequest)
			return
		}
	}

	if _, ok := fieldsToUpdate["bike_id"]; ok {
		isBikeAvailable := h.RentalRepo.IsBikeAvailable(*updateRentalReq.BikeID)
		if !isBikeAvailable {
			http.Error(w, "Bike is already rented by other user", http.StatusBadRequest)
			return
		}
	}

	result, err := h.RentalRepo.UpdateRental(rentalID, fieldsToUpdate)
	if err != nil {
		http.Error(w, "Error updating rental", http.StatusBadRequest)
		return
	}
	updateRentalResp := models.UpdateRentalResponse{
		ID:      result,
		Message: "Rental updated successfully",
	}
	helpers.WriteJSON(w, http.StatusOK, updateRentalResp)

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

// GetRentalHistoryByUserID retrieves the rental history of a user
func (h *handler) GetRentalHistoryByUserID(w http.ResponseWriter, req *http.Request) {
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
	pageID := req.Context().Value(middlewares.PageIDKey)

	rentals, err := h.RentalRepo.GetRentalHistoryByUserID(userId, pageID.(int64))
	if err != nil {
		log.Printf("Error getting rental history: %v", err)
		http.Error(w, "Error getting rental history", http.StatusBadRequest)
		return
	}
	helpers.WriteJSON(w, http.StatusOK, rentals)
}

// getFieldsToUpdate compares the fields of the update request with the rental and returns the fields to update as map
func getFieldsToUpdate(updateRentalReq *models.UpdateRentalRequest, rental *models.Rental) (map[string]interface{}, error) {
	if updateRentalReq == nil || rental == nil {
		return nil, fmt.Errorf("no fields to update")
	}
	fieldsToUpdate := make(map[string]interface{})
	if updateRentalReq.BikeID != nil && *updateRentalReq.BikeID != rental.BikeID {
		fieldsToUpdate["bike_id"] = updateRentalReq.BikeID
	}
	if updateRentalReq.UserID != nil && *updateRentalReq.UserID != rental.UserID {
		fieldsToUpdate["user_id"] = updateRentalReq.UserID
	}
	if updateRentalReq.StartLatitude != nil && *updateRentalReq.StartLatitude != rental.StartLatitude {
		fieldsToUpdate["start_latitude"] = updateRentalReq.StartLatitude
	}
	if updateRentalReq.StartLongitude != nil && *updateRentalReq.StartLongitude != rental.StartLongitude {
		fieldsToUpdate["start_longitude"] = updateRentalReq.StartLongitude
	}

	if len(fieldsToUpdate) == 0 {
		return nil, fmt.Errorf("no fields to update")
	}
	return fieldsToUpdate, nil
}
