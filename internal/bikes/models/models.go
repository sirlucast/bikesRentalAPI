package models

import (
	"time"
)

// Bike contains the information of a bike
type Bike struct {
	ID             int64     `json:"id,omitempty"`
	IsAvailable    bool      `json:"is_available,omitempty"`
	Latitude       float64   `json:"latitude,omitempty"`
	Longitude      float64   `json:"longitude,omitempty"`
	PricePerMinute float64   `json:"price_per_minute,omitempty"`
	CreatedAt      time.Time `json:"created_at,omitempty"`
	UpdatedAt      time.Time `json:"updated_at,omitempty"`
} // @name Bike

// BikeList contains a list of bikes
type BikeList struct {
	// The list of bikes
	Items []*Bike `json:"items"`
	// The id to query the next page
	NextPageID int64 `json:"next_page_id,omitempty" example:"10"`
} // @name BikeList
