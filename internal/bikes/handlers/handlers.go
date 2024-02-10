package handlers

import (
	"bikesRentalAPI/internal/bikes/models"
	"bikesRentalAPI/internal/bikes/repository"
	"bikesRentalAPI/internal/helpers"
	"bikesRentalAPI/internal/middlewares"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
)

type Handler interface {
	AddBike(w http.ResponseWriter, req *http.Request)
	UpdateBike(w http.ResponseWriter, req *http.Request)
	GetBikeByID(w http.ResponseWriter, req *http.Request)
	ListAllBikes(w http.ResponseWriter, req *http.Request)
	ListAvailableBikes(w http.ResponseWriter, req *http.Request)
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

// ListAvailableBikes for users returns a list of available bikes with next page ID for pagination
func (h *handler) ListAvailableBikes(w http.ResponseWriter, r *http.Request) {
	pageID := r.Context().Value(middlewares.PageIDKey)
	bikes, err := h.BikeRepo.ListAvailableBikes(pageID.(int64))
	if err != nil {
		log.Printf("Error getting available bikes: %v", err)
		http.Error(w, "Error getting available bikes", http.StatusInternalServerError)
		return
	}

	bikesListResponse := models.BikeList{
		Items:      bikes.Items,
		NextPageID: bikes.NextPageID,
	}

	helpers.WriteJSON(w, http.StatusOK, bikesListResponse)

}

// -------------------------------
// ------ Admin access only ------
// -------------------------------

// AddBike creates a new bike in the database
func (h *handler) AddBike(w http.ResponseWriter, req *http.Request) {
	// Parse body from request
	body, err := helpers.ParseBody(req.Body)
	if err != nil {
		http.Error(w, "Error parsing body request", http.StatusBadRequest)
		return
	}
	var newBike models.CreateUpdateBikeRequest
	err = json.Unmarshal(body, &newBike)
	if err != nil {
		http.Error(w, "Error unmarshalling body request", http.StatusBadRequest)
		return
	}
	// Validate bike input
	err = h.validator.Struct(newBike)
	if err != nil {
		errors := err.(validator.ValidationErrors)
		http.Error(w, fmt.Sprintf("Validation errors: %s", errors), http.StatusBadRequest)
		return
	}
	*newBike.IsAvailable = true // New bikes are available by default
	id, err := h.BikeRepo.CreateBike(newBike)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error creating user: %v", err), http.StatusInternalServerError)
		return
	}
	updateBikeRespnse := models.CreateUpdateBikeResponse{
		ID:      id,
		Message: "Bike created successfully",
	}
	helpers.WriteJSON(w, http.StatusCreated, updateBikeRespnse)
}

// UpdateBike updates a bike in the database
// the URL parameters 'bike_id' passed through as the request
func (h *handler) UpdateBike(w http.ResponseWriter, req *http.Request) {
	bikeIDStr := chi.URLParam(req, "bike_id")
	bikeID, err := strconv.ParseInt(bikeIDStr, 10, 64)
	if err != nil {
		http.Error(w, fmt.Sprintf("couldn't read %s: %v", bikeIDStr, err), http.StatusBadRequest)
		return
	}

	// Parse body from request
	body, err := helpers.ParseBody(req.Body)
	if err != nil {
		http.Error(w, "Error parsing body request", http.StatusBadRequest)
		return
	}
	var updateBikeReq *models.CreateUpdateBikeRequest
	err = json.Unmarshal(body, &updateBikeReq)
	if err != nil {
		http.Error(w, "Error unmarshalling body request", http.StatusBadRequest)
		return
	}
	// Validate user input
	err = h.validator.Struct(updateBikeReq)
	if err != nil {
		errors := err.(validator.ValidationErrors)
		http.Error(w, fmt.Sprintf("Validation errors: %s", errors), http.StatusBadRequest)
		return
	}

	bike, err := h.BikeRepo.GetBikeByID(bikeID)
	if err != nil {
		log.Printf("Error getting bike: %v", err)
		http.Error(w, "Error getting bike", http.StatusInternalServerError)
		return
	}

	// Compare attributes to identify modifications
	fieldsToUpdate, err := getFieldsToUpdate(updateBikeReq, bike)
	if err != nil {
		http.Error(w, "No fields to update", http.StatusBadRequest)
		return
	}

	// Update user
	result, err := h.BikeRepo.UpdateBike(bikeID, fieldsToUpdate)
	if err != nil {
		http.Error(w, "Error updating user", http.StatusInternalServerError)
		return
	}
	updateUserResp := models.CreateUpdateBikeResponse{
		ID:      result,
		Message: "Bike updated successfully",
	}
	helpers.WriteJSON(w, http.StatusOK, updateUserResp)
}

// getFieldsToUpdate compares the fields of the update request with the bike and returns the fields to update as map
func getFieldsToUpdate(updateBikeReq *models.CreateUpdateBikeRequest, bike *models.Bike) (map[string]interface{}, error) {
	if updateBikeReq == nil || bike == nil {
		return nil, fmt.Errorf("no fields to update")
	}
	fieldsToUpdate := make(map[string]interface{})
	if updateBikeReq.IsAvailable != nil && *updateBikeReq.IsAvailable != bike.IsAvailable {
		fieldsToUpdate["is_available"] = updateBikeReq.IsAvailable
	}
	if updateBikeReq.Latitude != nil && *updateBikeReq.Latitude != bike.Latitude {
		fieldsToUpdate["latitude"] = updateBikeReq.Latitude
	}
	if updateBikeReq.Longitude != nil && *updateBikeReq.Longitude != bike.Longitude {
		fieldsToUpdate["longitude"] = updateBikeReq.Longitude
	}
	if updateBikeReq.PricePerMinute != nil && *updateBikeReq.PricePerMinute != bike.PricePerMinute {
		fieldsToUpdate["price_per_minute"] = updateBikeReq.PricePerMinute

	}

	if len(fieldsToUpdate) == 0 {
		return nil, fmt.Errorf("no fields to update")
	}
	return fieldsToUpdate, nil
}

// ListAllBikes retrieves all bikes from the database
func (h *handler) ListAllBikes(w http.ResponseWriter, r *http.Request) {
	pageID := r.Context().Value(middlewares.PageIDKey)
	bikes, err := h.BikeRepo.ListAllBikes(pageID.(int64))
	if err != nil {
		log.Printf("Error getting available bikes: %v", err)
		http.Error(w, "Error getting available bikes", http.StatusInternalServerError)
		return
	}

	bikesListResponse := models.BikeList{
		Items:      bikes.Items,
		NextPageID: bikes.NextPageID,
	}
	helpers.WriteJSON(w, http.StatusOK, bikesListResponse)
}

// GetBikeByID retrieves a bike from the database
func (h *handler) GetBikeByID(w http.ResponseWriter, r *http.Request) {
	bikeIDStr := chi.URLParam(r, "bike_id")
	bikeID, err := strconv.ParseInt(bikeIDStr, 10, 64)
	if err != nil {
		http.Error(w, fmt.Sprintf("couldn't read %s: %v", bikeIDStr, err), http.StatusBadRequest)
		return
	}
	bike, err := h.BikeRepo.GetBikeByID(bikeID)
	if err != nil {
		log.Printf("Error getting bike: %v", err)
		http.Error(w, "Error getting bike", http.StatusInternalServerError)
		return
	}
	helpers.WriteJSON(w, http.StatusOK, bike)
}
