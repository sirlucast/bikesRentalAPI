package main

import (
	"bikesRentalAPI/internal/database"
	"bikesRentalAPI/internal/router"
	"bikesRentalAPI/internal/server"
	"log"

	"github.com/golang-migrate/migrate"
)

func main() {

	// Create a new database service
	dbService := database.New()

	err := dbService.Start()
	if err != nil {
		log.Fatalf("failed to start database: %v", err)
	}
	defer func() {
		if err := dbService.Close(); err != nil {
			log.Fatalf("Failed to close database connection: %v", err)
		}
	}()

	// Check DB health
	err = dbService.Health()
	if err != nil {
		log.Fatalf("failed to check database health: %v", err)
	}

	// Migrate the database
	err = dbService.Migrate()
	if err != nil && err != migrate.ErrNoChange {
		log.Fatalf("failed to migrate database: %v", err)
	}
	// Create a new router service
	routerService := router.New()

	server, err := server.NewServerBuilder().
		WithRouter(routerService).
		Build()
	if err != nil {
		log.Fatalf("failed to build server: %v", err)
	}

	log.Printf("Server running on port %s", server.Addr)
	err = server.ListenAndServe()
	if err != nil {
		log.Fatalf("cannot start server: %s", err)
	}
}
