package repository

import (
	bikesrepository "bikesRentalAPI/internal/bikes/repository"
	"bikesRentalAPI/internal/database"
	"bikesRentalAPI/internal/helpers"
	"bikesRentalAPI/internal/rentals/models"
	usersrepository "bikesRentalAPI/internal/users/repository"
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"
)

const (
	// pageSize is the number of items to return in a page, 10 as default.
	pageSize = 10
)

type RentalRepository interface {
	IsBikeAvailable(bikeID int64) bool
	IsUserRentingBike(userID int64) bool
	StartRental(userID int64, startReq *models.StartBikeRentalRequest) (*models.StartRentalResponse, error)
	EndRental(userID int64, endReq *models.StopBikeRentalRequest) (*models.StopRentalResponse, error)
	GetOngoingRental(userID int64) (*models.Rental, error)
}

type rentalRepository struct {
	db       database.Database
	userRepo usersrepository.UserRepository
	bikeRepo bikesrepository.BikeRepository
}

// New initializes a new empty rental repository
func New(db database.Database,
	userRepo usersrepository.UserRepository,
	bikeRepo bikesrepository.BikeRepository,
) RentalRepository {
	return &rentalRepository{
		db:       db,
		userRepo: userRepo,
		bikeRepo: bikeRepo,
	}
}

func (r *rentalRepository) IsBikeAvailable(bikeID int64) bool {
	isAvailable, err := r.bikeRepo.IsBikeAvailable(bikeID)
	if err != nil {
		log.Printf("Error checking if bike is available: %v", err)
		return false
	}
	return isAvailable
}

func (r *rentalRepository) IsUserRentingBike(userID int64) bool {
	var count int
	query := "SELECT COUNT(*) FROM rentals WHERE user_id = ? AND end_time IS NULL"
	err := r.db.QueryRow(query, userID).Scan(&count)
	if err != nil {
		log.Printf("Error checking if user is renting a bike: %v", err)
		return false
	}
	return count > 0
}

func (r *rentalRepository) StartRental(userID int64, startReq *models.StartBikeRentalRequest) (*models.StartRentalResponse, error) {
	if startReq == nil {
		return nil, fmt.Errorf("startReq request is nil")
	}
	now := time.Now().UTC()
	initialCost := 0.0
	query := "INSERT INTO rentals (user_id, bike_id, start_time, start_latitude, start_longitude, cost) VALUES (?, ?, ?, ?, ?, ?)"
	result, err := r.db.Exec(query, userID, startReq.BikeID, now, startReq.Latitude, startReq.Longitude, initialCost)
	if err != nil {
		return nil, fmt.Errorf("failed to insert rental: %v", err)
	}
	err = r.bikeRepo.SetBikeAvailability(startReq.BikeID, false)
	if err != nil {
		return nil, fmt.Errorf("failed to set bike availability: %v", err)
	}
	id, err := result.LastInsertId()
	if err != nil {
		return nil, fmt.Errorf("failed to get last insert id: %v", err)
	}
	return &models.StartRentalResponse{
		ID:        id,
		StartTime: now,
		Latitude:  startReq.Latitude,
		Longitude: startReq.Longitude,
	}, nil

}

func (r *rentalRepository) GetOngoingRental(userID int64) (*models.Rental, error) {
	var rental models.Rental
	query := "SELECT id, user_id, bike_id, start_time, start_latitude, start_longitude, cost FROM rentals WHERE user_id = ? AND end_time IS NULL ORDER BY start_time DESC LIMIT 1"
	err := r.db.QueryRow(query, userID).Scan(&rental.ID, &rental.UserID, &rental.BikeID, &rental.StartTime, &rental.StartLatitude, &rental.StartLongitude, &rental.Cost)
	if err != nil {
		return nil, err
	}
	return &rental, nil
}

func (r *rentalRepository) EndRental(userID int64, endReq *models.StopBikeRentalRequest) (*models.StopRentalResponse, error) {
	if endReq == nil {
		return nil, fmt.Errorf("endReq request is nil")
	}
	rental, err := r.GetOngoingRental(userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get ongoing rental: %v", err)
	}

	bikeCostPerMin, err := r.bikeRepo.GetBikeCostPerMinute(rental.BikeID)
	if err != nil {
		return nil, fmt.Errorf("failed to get bike cost per minute: %v", err)
	}

	now := time.Now().UTC()
	finalLat, finalLon := helpers.GetRandomLatLon(rental.StartLatitude, rental.StartLongitude)
	duration := now.Sub(rental.StartTime)
	cost := calculateRentalCost(bikeCostPerMin, duration)

	r.db.Transaction(context.Background(), func(tx *sql.Tx) error {

		query := "UPDATE rentals SET end_time = ?, end_latitude = ?, end_longitude = ?, duration_minutes = ?, cost = ? WHERE id = ?"
		_, err = r.db.Exec(query, now, finalLat, finalLon, duration.Minutes(), cost, rental.ID)
		if err != nil {
			return fmt.Errorf("failed to update rental: %v", err)
		}
		err = r.bikeRepo.SetBikeAvailability(rental.BikeID, true)
		if err != nil {
			return fmt.Errorf("failed to set bike availability: %v", err)
		}
		return nil
	})

	return &models.StopRentalResponse{
		BikeID:          rental.BikeID,
		EndTime:         now,
		Latitude:        finalLat,
		Longitude:       finalLon,
		Cost:            cost,
		DurationMinutes: int(duration.Round(time.Minute).Minutes()),
	}, nil
}

func calculateRentalCost(costPerMin float64, duration time.Duration) float64 {
	return costPerMin * duration.Minutes()
}
