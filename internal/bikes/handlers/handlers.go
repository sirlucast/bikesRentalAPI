package handlers

import "net/http"

// ListAvailableBikes ...
func ListAvailableBikes(w http.ResponseWriter, r *http.Request) {
	// TODO Implement logic to list available bikes
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
func AddBike(w http.ResponseWriter, r *http.Request) {
	// TODO Implement logic to add a new bike
}

// UpdateBike ...
func UpdateBike(w http.ResponseWriter, r *http.Request) {
	// TODO Implement logic to update bike details
}

// DeleteBike ...
func ListBikes(w http.ResponseWriter, r *http.Request) {
	// TODO Implement logic to list all bikes
}
