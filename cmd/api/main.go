package main

import (
	"bikesRentalAPI/internal/database"
	"bikesRentalAPI/internal/router"
	"bikesRentalAPI/internal/server"
	userhandler "bikesRentalAPI/internal/users/handlers"
	usermodels "bikesRentalAPI/internal/users/models"
	userrepository "bikesRentalAPI/internal/users/repository"
	"log"

	_ "github.com/joho/godotenv/autoload"
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
	if err != nil {
		log.Fatalf("failed to migrate database: %v", err)
	}

	// Seeds the database
	seeder := database.NewSeeder(dbService)
	err = seeder.Seed(usermodels.User{})
	if err != nil {
		log.Fatalf("failed to seed database: %v", err)
	}

	// Initialize the repositories
	userRepository := userrepository.New(dbService)
	userHandler := userhandler.New(userRepository)

	// Create a new router service and register routes
	routerService := router.New()
	handler := routerService.RegisterRoutes(userHandler)

	server, err := server.NewServerBuilder().
		WithHanlder(handler).
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
