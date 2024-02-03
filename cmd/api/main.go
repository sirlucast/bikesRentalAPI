package main

import (
	"bikesRentalAPI/internal/server"
	"fmt"
	"log"
)

func main() {

	server := server.NewServer()
	log.Printf("Server running on port %s", server.Addr)
	err := server.ListenAndServe()
	if err != nil {
		log.Fatalf(fmt.Sprintf("cannot start server: %s", err))
	}
}
