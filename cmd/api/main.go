package main

import (
	"bikesRentalAPI/internal/database"
	"bikesRentalAPI/internal/router"
	"bikesRentalAPI/internal/server"
	"log"
)

func main() {

	// Create a new database service
	dbService, err := database.New()
	if err != nil {
		log.Fatalf("failed to start database: %v", err)
	}

	// Create a new router service
	routerService := router.New()

	dbService.Health()

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
