package handlers

import "net/http"

// Admin access only

// ListRentals ...
func ListRentals(w http.ResponseWriter, r *http.Request) {
	// TODO Implement logic to list all rentals
}

// GetRentalDetails ...
func GetRentalDetails(w http.ResponseWriter, r *http.Request) {
	// TODO Implement logic to get rental details
}

// UpdateRentalDetails ...
func UpdateRentalDetails(w http.ResponseWriter, r *http.Request) {
	// TODO Implement logic to update rental details
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
	// 1. Get pagination details from the request context
	//pageID := r.Context().Value(middlewares.PageIDKey)
	// 2. Get rental history from the database
}
