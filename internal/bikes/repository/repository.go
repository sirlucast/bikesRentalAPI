package repository

import (
	"bikesRentalAPI/internal/bikes/models"
	"bikesRentalAPI/internal/database"
)

type BikeRepository interface {
	ListAvailableBikes() ([]models.Bike, error)
}

type bikeRepository struct {
	db database.Database
}

// New initializes a new empty bike repository
func New(db database.Database) BikeRepository {
	return &bikeRepository{db}
}

func (r *bikeRepository) ListAvailableBikes() ([]models.Bike, error) {
	query := "SELECT * FROM bikes WHERE available = ?"
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
