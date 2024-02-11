package handlers

import (
	"bikesRentalAPI/internal/rentals/repository"
	"net/http"

	"github.com/go-playground/validator/v10"
)

// Handler is the interface for rental handlers
type Handler interface {
	GetRentalHistoryByUserID(w http.ResponseWriter, r *http.Request) // Get rental history by user ID
	GetRentalList(w http.ResponseWriter, r *http.Request)            // Get rental list
	GetRentalDetails(w http.ResponseWriter, r *http.Request)         // Get rental details
	UpdateRentalDetails(w http.ResponseWriter, r *http.Request)      // Update rental details
	StartBikeRental(w http.ResponseWriter, r *http.Request)          // Start bike rental
	EndBikeRental(w http.ResponseWriter, r *http.Request)            // End bike rental
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
func (h *handler) GetRentalList(w http.ResponseWriter, r *http.Request) {
	// TODO Implement logic to list all rentals
}

// GetRentalDetails ...
func (h *handler) GetRentalDetails(w http.ResponseWriter, r *http.Request) {
	// TODO Implement logic to get rental details
}

// UpdateRentalDetails ...
func (h *handler) UpdateRentalDetails(w http.ResponseWriter, r *http.Request) {
	// TODO Implement logic to update rental details
}

// StartBikeRental ...
func (h *handler) StartBikeRental(w http.ResponseWriter, r *http.Request) {
	// TODO Implement logic to start bike rental
}

// EndBikeRental ...
func (h *handler) EndBikeRental(w http.ResponseWriter, r *http.Request) {
	// TODO Implement logic to end bike rental
}

// GetRentalHistoryByUserID ...
func (h *handler) GetRentalHistoryByUserID(w http.ResponseWriter, r *http.Request) {
	// TODO Implement logic to get rental history
	// 1. Get pagination details from the request context
	//pageID := r.Context().Value(middlewares.PageIDKey)
	// 2. Get rental history from the database
}
