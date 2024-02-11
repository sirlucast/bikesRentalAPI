package models

import "time"

// Rental model represents a rental operation of a bike
type Rental struct {
	// The id of the rental
	ID int64 `json:"id"`
	// The id of the user
	UserID int64 `json:"user_id"`
	// The id of the bike
	BikeID int64 `json:"bike_id"`
	// The start time of the rental
	StartTime *time.Time `json:"start_time"`
	//  The end time of the rental
	EndTime *time.Time `json:"end_time"`
	// The latitude of the start location
	StartLatitude float64 `json:"start_latitude"`
	// The longitude of the start location
	StartLongitude float64 `json:"start_longitude"`
	// The latitude of the end location
	EndLatitude *float64 `json:"end_latitude"`
	// The longitude of the end location
	EndLongitude *float64 `json:"end_longitude"`
	// The creation time of the rental
	CreatedAt *time.Time `json:"created_at"`
	// The last update time of the rental
	UpdatedAt *time.Time `json:"updated_at"`
	// The duration of the rental in minutes
	DurationMinutes int64 `json:"duration_minutes"`
	// The cost of the rental
	Cost float64 `json:"cost"`
} // @name Rental

// StartBikeRentalRequest contains the request to start a rental
type StartBikeRentalRequest struct {
	BikeID    int64   `json:"bike_id" validate:"required,numeric"`
	Latitude  float64 `json:"latitude" validate:"required,latitude"`
	Longitude float64 `json:"longitude" validate:"required,longitude"`
} // @name StartBikeRentalRequest

// StopBikeRentalRequest contains the request to stop a rental
type StopBikeRentalRequest struct {
	RentalID int64 `json:"rental_id" validate:"required,numeric"`
} // @name StopBikeRentalRequest

// StartRentalResponse contains the response of starting a rental
type StartRentalResponse struct {
	ID        int64     `json:"id"`
	StartTime time.Time `json:"start_time"`
	Latitude  float64   `json:"latitude"`
	Longitude float64   `json:"longitude"`
} // @name StartRentalResponse

// StopRentalResponse contains the response of stopping a rental
type StopRentalResponse struct {
	BikeID          int64     `json:"bike_id"`
	EndTime         time.Time `json:"end_time"`
	Latitude        float64   `json:"latitude"`
	Longitude       float64   `json:"longitude"`
	Cost            float64   `json:"cost"`
	DurationMinutes int       `json:"duration"`
	Distance        float64   `json:"distance,omitempty"`
} // @name StopRentalResponse

// RentalList contains a list of rentals and the next page id
type RentalList struct {
	// The list of rentals
	Items []*Rental `json:"items"`
	// The id to query the next page
	NextPageID int64 `json:"next_page_id"`
} // @name RentalList

// UpdateRentalRequest contains the request to update a rental
type UpdateRentalRequest struct {
	RentalID       *int64     `json:"rental_id" validate:"required,numeric"`
	UserID         *int64     `json:"user_id" validate:"omitempty,required,numeric"`
	BikeID         *int64     `json:"bike_id" validate:"omitempty,required,numeric"`
	StartTime      *time.Time `json:"start_time" validate:"omitempty,required"`
	StartLatitude  *float64   `json:"start_latitude" validate:"omitempty,required,latitude"`
	StartLongitude *float64   `json:"start_longitude" validate:"omitempty,required,longitude"`
} // @name UpdateRentalRequest

// UpdateRentalResponse represents the response of updating a rental
type UpdateRentalResponse struct {
	ID      int64  `json:"id,omitempty"`
	Message string `json:"message,omitempty"`
} // @name UpdateRentalResponse
