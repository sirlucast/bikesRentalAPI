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
	"strings"
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
	GetRentalHistoryByUserID(userID int64, PageID int64) (*models.RentalList, error)
	GetRentalDetails(rentalID int64) (*models.Rental, error)
	UpdateRental(rentalID int64, fieldsToUpdate map[string]interface{}) (int64, error)
	ListAllRentals(pageID int64) (*models.RentalList, error)
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
	var id int64

	r.db.Transaction(context.Background(), func(tx *sql.Tx) error {
		query := "INSERT INTO rentals (user_id, bike_id, start_time, start_latitude, start_longitude, cost) VALUES (?, ?, ?, ?, ?, ?)"
		result, err := r.db.Exec(query, userID, startReq.BikeID, now, startReq.Latitude, startReq.Longitude, initialCost)
		if err != nil {
			return fmt.Errorf("failed to insert rental: %v", err)
		}
		err = r.bikeRepo.SetBikeAvailability(startReq.BikeID, false)
		if err != nil {
			return fmt.Errorf("failed to set bike availability: %v", err)
		}
		id, err = result.LastInsertId()
		if err != nil {
			return fmt.Errorf("failed to get last insert id: %v", err)
		}
		return nil
	})

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

	// For this example, we will generate random latitude and longitude for the end location.
	finalLat, finalLon := helpers.GetRandomLatLon(rental.StartLatitude, rental.StartLongitude)

	now := time.Now().UTC()
	duration := now.Sub(rental.StartTime.UTC())
	durationInMinutes := int(duration.Round(time.Minute).Minutes())
	cost := calculateRentalCost(bikeCostPerMin, duration)

	r.db.Transaction(context.Background(), func(tx *sql.Tx) error {

		query := "UPDATE rentals SET end_time = ?, end_latitude = ?, end_longitude = ?, duration_minutes = ?, cost = ? WHERE id = ?"
		_, err = r.db.Exec(query, now, finalLat, finalLon, durationInMinutes, cost, rental.ID)
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
		DurationMinutes: durationInMinutes,
	}, nil
}

// calculateRentalCost calculates the cost of a rental based on the cost per minute and the duration of the rental
func calculateRentalCost(costPerMin float64, duration time.Duration) float64 {
	return costPerMin * duration.Minutes()
}

// GetRentalHistoryByUserID returns a list of rentals for a user, paginated
func (r *rentalRepository) GetRentalHistoryByUserID(userID int64, PageID int64) (*models.RentalList, error) {

	query := "SELECT id, user_id, bike_id, start_time, end_time, start_latitude, start_longitude, end_latitude, end_longitude, duration_minutes, cost, created_at, updated_at FROM rentals WHERE user_id = ? AND id > ? ORDER BY id LIMIT ?"
	rows, err := r.db.Query(query, userID, PageID, pageSize)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	rentals := &models.RentalList{}
	rentalList := make([]*models.Rental, 0)

	for rows.Next() {
		var rental models.Rental
		if err := rows.Scan(&rental.ID,
			&rental.UserID,
			&rental.BikeID,
			&rental.StartTime,
			&rental.EndTime,
			&rental.StartLatitude,
			&rental.StartLongitude,
			&rental.EndLatitude,
			&rental.EndLongitude,
			&rental.DurationMinutes,
			&rental.Cost,
			&rental.CreatedAt,
			&rental.UpdatedAt,
		); err != nil {
			return nil, err
		}
		rentalList = append(rentalList, &rental)
	}
	if len(rentalList) == pageSize {
		rentals.NextPageID = rentalList[len(rentalList)-1].ID
	}
	rentals.Items = rentalList
	return rentals, nil
}

// GetRentalDetails returns the details of a rental
func (r *rentalRepository) GetRentalDetails(rentalID int64) (*models.Rental, error) {
	query := "SELECT id, user_id, bike_id, start_time, end_time, start_latitude, start_longitude, end_latitude, end_longitude, duration_minutes, cost, created_at, updated_at FROM rentals WHERE id = ?"
	row := r.db.QueryRow(query, rentalID)
	var rental models.Rental
	if err := row.Scan(&rental.ID,
		&rental.UserID,
		&rental.BikeID,
		&rental.StartTime,
		&rental.EndTime,
		&rental.StartLatitude,
		&rental.StartLongitude,
		&rental.EndLatitude,
		&rental.EndLongitude,
		&rental.DurationMinutes,
		&rental.Cost,
		&rental.CreatedAt,
		&rental.UpdatedAt,
	); err != nil {
		return nil, fmt.Errorf("failed to get rental details: %v", err)
	}
	return &rental, nil
}

// UpdateRental updates a rental in the database by id. Returns the id of the updated rental. If no fields are updated, returns 0 and an error
func (r *rentalRepository) UpdateRental(rentalID int64, fieldsToUpdate map[string]interface{}) (int64, error) {
	var setFields []string
	var args []interface{}

	for field, value := range fieldsToUpdate {
		setFields = append(setFields, fmt.Sprintf("%s = ?", field))
		args = append(args, value)
	}
	args = append(args, rentalID)

	query := fmt.Sprintf("UPDATE rentals SET %s WHERE id = ?", strings.Join(setFields, ", "))
	stmt, err := r.db.Prepare(query)
	if err != nil {
		return 0, fmt.Errorf("failed to prepare update statement: %v", err)
	}
	defer stmt.Close()

	result, err := stmt.Exec(args...)
	if err != nil {
		return 0, fmt.Errorf("failed to execute update statement: %v", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("failed to get last insert id: %v", err)
	}
	return id, nil
}

func (r *rentalRepository) ListAllRentals(pageID int64) (*models.RentalList, error) {
	query := "SELECT id, user_id, bike_id, start_time, end_time, start_latitude, start_longitude, end_latitude, end_longitude, duration_minutes, cost, created_at, updated_at FROM rentals WHERE id > ? ORDER BY id LIMIT ?"
	rows, err := r.db.Query(query, pageID, pageSize)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	rentals := &models.RentalList{}
	rentalList := make([]*models.Rental, 0)

	for rows.Next() {
		var rental models.Rental
		if err := rows.Scan(&rental.ID,
			&rental.UserID,
			&rental.BikeID,
			&rental.StartTime,
			&rental.EndTime,
			&rental.StartLatitude,
			&rental.StartLongitude,
			&rental.EndLatitude,
			&rental.EndLongitude,
			&rental.DurationMinutes,
			&rental.Cost,
			&rental.CreatedAt,
			&rental.UpdatedAt,
		); err != nil {
			return nil, err
		}
		rentalList = append(rentalList, &rental)
	}
	if len(rentalList) == pageSize {
		rentals.NextPageID = rentalList[len(rentalList)-1].ID
	}
	rentals.Items = rentalList
	return rentals, nil
}
