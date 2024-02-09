package repository

import (
	"bikesRentalAPI/internal/bikes/models"
	"bikesRentalAPI/internal/database"
	"fmt"
	"strings"
)

type BikeRepository interface {
	ListAvailableBikes() ([]models.Bike, error)
	ListAllBikes() ([]models.Bike, error)
	UpdateBike(bikeID int, fieldsToUpdate map[string]interface{}) (int64, error)
	CreateBike(bike models.Bike) (int64, error)
}

type bikeRepository struct {
	db database.Database
}

// New initializes a new empty bike repository
func New(db database.Database) BikeRepository {
	return &bikeRepository{db}
}

// ListAvailableBikes retrieves all available bikes from the database
func (r *bikeRepository) ListAvailableBikes() ([]models.Bike, error) {
	query := "SELECT id, is_available, price_per_minute FROM bikes WHERE is_available = ?"
	rows, err := r.db.Query(query, true)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var bikes []models.Bike
	for rows.Next() {
		var bike models.Bike
		if err := rows.Scan(&bike.ID, &bike.IsAvailable, &bike.PricePerMinute); err != nil {
			return nil, err
		}
		bikes = append(bikes, bike)
	}
	return bikes, nil
}

// ListAllBikes retrieves all bikes from the database
func (r *bikeRepository) ListAllBikes() ([]models.Bike, error) {
	query := "SELECT * FROM bikes"
	rows, err := r.db.Query(query, true)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var bikes []models.Bike
	for rows.Next() {
		var bike models.Bike
		if err := rows.Scan(&bike.ID, &bike.IsAvailable, &bike.PricePerMinute, &bike.Latitude, &bike.Longitude, &bike.UpdatedAt); err != nil {
			return nil, err
		}
		bikes = append(bikes, bike)
	}
	return bikes, nil
}

// UpdateBike updates a bike in the database
func (r *bikeRepository) UpdateBike(bikeID int, fieldsToUpdate map[string]interface{}) (int64, error) {
	var setFields []string
	var args []interface{}

	for field, value := range fieldsToUpdate {
		setFields = append(setFields, fmt.Sprintf("%s = ?", field))
		args = append(args, value)
	}
	args = append(args, bikeID)

	query := fmt.Sprintf("UPDATE bikes SET %s WHERE id = ?", strings.Join(setFields, ", "))
	stmt, err := r.db.Prepare(query)
	if err != nil {
		return 0, fmt.Errorf("failed to prepare update statement: %v", err)
	}
	defer stmt.Close()

	result, err := r.db.Exec(query, args...)
	if err != nil {
		return 0, err
	}
	lastInsertedID, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}
	return lastInsertedID, nil
}

// CreateBike creates a bike in the database
func (r *bikeRepository) CreateBike(bike models.Bike) (int64, error) {
	query := "INSERT INTO bikes (is_available, price_per_minute, latitude, longitude) VALUES (?, ?, ?, ?)"
	result, err := r.db.Exec(query, bike.IsAvailable, bike.PricePerMinute, bike.Latitude, bike.Longitude)
	if err != nil {
		return 0, fmt.Errorf("failed to insert bike: %v", err)
	}
	id, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("failed to get last insert id: %v", err)
	}
	return id, nil
}
