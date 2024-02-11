package models

import "time"

// Rental model
type Rental struct {
	ID              int64     `json:"id"`
	UserID          int64     `json:"user_id"`
	BikeID          int64     `json:"bike_id"`
	StartTime       time.Time `json:"start_time"`
	EndTime         time.Time `json:"end_time"`
	StartLatitude   float64   `json:"start_latitude"`
	StartLongitude  float64   `json:"start_longitude"`
	EndLatitude     float64   `json:"end_latitude"`
	EndLongitude    float64   `json:"end_longitude"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
	DurationMinutes int64     `json:"duration_minutes"`
	Cost            float64   `json:"cost"`
}

type StartBikeRentalRequest struct {
	BikeID    int64   `json:"bike_id" validate:"required,numeric"`
	Latitude  float64 `json:"latitude" validate:"required,latitude"`
	Longitude float64 `json:"longitude" validate:"required,longitude"`
}

type StopBikeRentalRequest struct {
	RentalID int64 `json:"rental_id" validate:"required,numeric"`
}

type StartRentalResponse struct {
	ID        int64     `json:"id"`
	StartTime time.Time `json:"start_time"`
	Latitude  float64   `json:"latitude"`
	Longitude float64   `json:"longitude"`
}

type StopRentalResponse struct {
	BikeID          int64     `json:"bike_id"`
	EndTime         time.Time `json:"end_time"`
	Latitude        float64   `json:"latitude"`
	Longitude       float64   `json:"longitude"`
	Cost            float64   `json:"cost"`
	DurationMinutes int       `json:"duration"`
	Distance        float64   `json:"distance,omitempty"`
}
