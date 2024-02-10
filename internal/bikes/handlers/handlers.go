package handlers

import (
	"bikesRentalAPI/internal/bikes/repository"
	"bikesRentalAPI/internal/helpers"
	"bikesRentalAPI/internal/middlewares"
	"log"
	"net/http"

	"github.com/go-playground/validator/v10"
)

type Handler interface {
	ListAvailableBikes(w http.ResponseWriter, req *http.Request)
	AddBike(w http.ResponseWriter, req *http.Request)
	UpdateBike(w http.ResponseWriter, req *http.Request)
	ListAllBikes(w http.ResponseWriter, req *http.Request)
}

type handler struct {
	BikeRepo  repository.BikeRepository
	validator *validator.Validate
}

// New returns a new user handler
func New(BikeRepository repository.BikeRepository) Handler {
	validator := validator.New(validator.WithRequiredStructEnabled())
	handler := &handler{
		BikeRepo:  BikeRepository,
		validator: validator,
	}
	return handler
}

// ListAvailableBikes for users
func (h *handler) ListAvailableBikes(w http.ResponseWriter, r *http.Request) {
	pageID := r.Context().Value(middlewares.PageIDKey)
	bikes, err := h.BikeRepo.ListAvailableBikes(pageID.(int64))
	if err != nil {
		log.Printf("Error getting available bikes: %v", err)
		http.Error(w, "Error getting available bikes", http.StatusInternalServerError)
		return
	}

	helpers.WriteJSON(w, http.StatusOK, bikes)

}

// StartBikeRental ...
func StartBikeRental(w http.ResponseWriter, r *http.Request) {
	// TODO Implement logic to start bike rental
}

// EndBikeRental ...
func EndBikeRental(w http.ResponseWriter, r *http.Request) {
	// TODO Implement logic to end bike rental
}

// GetRentalHistory ...
func GetRentalHistory(w http.ResponseWriter, r *http.Request) {
	// TODO Implement logic to get rental history
}

// Admin access only

// AddBike ...
func (h *handler) AddBike(w http.ResponseWriter, r *http.Request) {
	// TODO Implement logic to add a new bike
}

// UpdateBike ...
func (h *handler) UpdateBike(w http.ResponseWriter, r *http.Request) {
	// TODO Implement logic to update bike details
}

// DeleteBike ...
func (h *handler) ListAllBikes(w http.ResponseWriter, r *http.Request) {
	// TODO Implement logic to list all bikes
}
