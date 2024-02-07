package models

import "time"

// Rental model
type Rental struct {
	ID              int64     `json:"id"`
	UserID          int64     `json:"user_id"`
	BikeID          int64     `json:"bike_id"`
	Status          string    `json:"status"`
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
