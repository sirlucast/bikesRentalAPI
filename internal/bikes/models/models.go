package models

import "time"

type Bike struct {
	ID             int64     `json:"id,omitempty"`
	IsAvailable    bool      `json:"is_available,omitempty"`
	Latitude       float64   `json:"latitude,omitempty"`
	Longitude      float64   `json:"longitude,omitempty"`
	CreatedAt      time.Time `json:"created_at,omitempty"`
	UpdatedAt      time.Time `json:"updated_at,omitempty"`
	PricePerMinute float64   `json:"price_per_minute,omitempty"`
}
