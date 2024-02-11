package repository

import "bikesRentalAPI/internal/database"

const (
	// pageSize is the number of items to return in a page, 10 as default.
	pageSize = 10
)

type RentalRepository interface {
}

type rentalRepository struct {
	db database.Database
}

// New initializes a new empty rental repository
func New(db database.Database) RentalRepository {
	return &rentalRepository{db}
}
