package models

import "time"

type Bike struct {
	ID             int64     `json:"id"`
	IsAvailable    bool      `json:"is_available"`
	Latitude       float64   `json:"latitude"`
	Longitude      float64   `json:"longitude"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
	PricePerMinute float64   `json:"price_per_minute"`
}
