package repository

import (
	"bikesRentalAPI/internal/bikes/models"
	"bikesRentalAPI/internal/database"
	"fmt"
	"log"
	"strings"
)

const (
	// pageSize is the number of items to return in a page, 10 as default.
	pageSize = 10
)

type BikeRepository interface {
	ListAvailableBikes(PageID int64) (*models.BikeList, error)
	ListAllBikes(PageID int64) (*models.BikeList, error)
	UpdateBike(bikeID int64, fieldsToUpdate map[string]interface{}) (int64, error)
	CreateBike(bike models.CreateUpdateBikeRequest) (int64, error)
	GetBikeByID(bikeID int64) (*models.Bike, error)
}

type bikeRepository struct {
	db database.Database
}

// New initializes a new empty bike repository
func New(db database.Database) BikeRepository {
	return &bikeRepository{db}
}

// GetBikeByID retrieves a bike from the database by its id
func (r *bikeRepository) GetBikeByID(bikeID int64) (*models.Bike, error) {
	query := "SELECT id, is_available, price_per_minute, latitude, longitude, created_at, updated_at FROM bikes WHERE id = ?"
	row := r.db.QueryRow(query, bikeID)
	var bike models.Bike
	if err := row.Scan(&bike.ID, &bike.IsAvailable, &bike.PricePerMinute, &bike.Latitude, &bike.Longitude, &bike.CreatedAt, &bike.UpdatedAt); err != nil {
		return nil, err
	}
	return &bike, nil
}

// ListAvailableBikes retrieves all available bikes from the database
func (r *bikeRepository) ListAvailableBikes(PageID int64) (*models.BikeList, error) {
	query := "SELECT id, is_available, price_per_minute FROM bikes WHERE is_available = ? AND id > ? ORDER BY id LIMIT ?"
	rows, err := r.db.Query(query, true, PageID, pageSize)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	bikes := &models.BikeList{}
	bikeList := make([]*models.Bike, 0)
	for rows.Next() {
		var bike models.Bike
		if err := rows.Scan(&bike.ID, &bike.IsAvailable, &bike.PricePerMinute); err != nil {
			return nil, err
		}

		bikeList = append(bikeList, &bike)
	}
	if len(bikeList) == pageSize {
		bikes.NextPageID = bikeList[len(bikeList)-1].ID
	}
	bikes.Items = bikeList
	log.Printf("bikes: %v", bikes)
	return bikes, nil
}

// ListAllBikes retrieves all bikes from the database
func (r *bikeRepository) ListAllBikes(PageID int64) (*models.BikeList, error) {
	query := "SELECT id, is_available, price_per_minute, latitude, longitude, created_at, updated_at FROM bikes WHERE id > ? ORDER BY id LIMIT ?"
	rows, err := r.db.Query(query, PageID, pageSize)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	bikes := &models.BikeList{}
	bikeList := make([]*models.Bike, 0)
	for rows.Next() {
		var bike models.Bike
		if err := rows.Scan(&bike.ID, &bike.IsAvailable, &bike.PricePerMinute, &bike.Latitude, &bike.Longitude, &bike.CreatedAt, &bike.UpdatedAt); err != nil {
			return nil, err
		}
		bikeList = append(bikeList, &bike)
	}
	if len(bikeList) == pageSize {
		bikes.NextPageID = bikeList[len(bikeList)-1].ID
	}
	bikes.Items = bikeList
	return bikes, nil
}

// UpdateBike updates a bike in the database
func (r *bikeRepository) UpdateBike(bikeID int64, fieldsToUpdate map[string]interface{}) (int64, error) {
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
func (r *bikeRepository) CreateBike(bike models.CreateUpdateBikeRequest) (int64, error) {
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
