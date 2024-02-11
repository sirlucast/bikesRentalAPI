package main

import (
	bikehandler "bikesRentalAPI/internal/bikes/handlers"
	bikerepository "bikesRentalAPI/internal/bikes/repository"
	"bikesRentalAPI/internal/database"
	rentalhanlder "bikesRentalAPI/internal/rentals/handlers"
	rentalrepository "bikesRentalAPI/internal/rentals/repository"
	"bikesRentalAPI/internal/router"
	"bikesRentalAPI/internal/server"
	userhandler "bikesRentalAPI/internal/users/handlers"
	userrepository "bikesRentalAPI/internal/users/repository"
	"log"

	_ "github.com/joho/godotenv/autoload"
)

func init() {
	log.SetFlags(log.LstdFlags | log.Llongfile)
}

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

	// - will move this to a better place
	// Seeds the database
	//seeder := database.NewSeeder(dbService)
	//err = seeder.SeedUser()
	//if err != nil {
	//	log.Fatalf("failed to seed database: %v", err)
	//}

	//bikesSeeder := database.NewSeeder(dbService)
	//err = bikesSeeder.SeedBikes()
	//if err != nil {
	//	log.Fatalf("failed to seed database: %v", err)
	//}
	// - end of seeds

	// Initialize the repositories and handlers
	userRepository := userrepository.New(dbService)
	userHandler := userhandler.New(userRepository)
	bikeRepository := bikerepository.New(dbService)
	bikeHandler := bikehandler.New(bikeRepository)

	rentalRepository := rentalrepository.New(dbService, userRepository, bikeRepository)
	rentalHanlder := rentalhanlder.New(rentalRepository)

	// Create a new router service and register routes
	routerService := router.New()
	handler := routerService.RegisterRoutes(userHandler, bikeHandler, rentalHanlder)

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
