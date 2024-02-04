package server

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"bikesRentalAPI/internal/router"

	_ "github.com/joho/godotenv/autoload"
)

// ServerConfig holds configuration parameters for the server
type ServerConfig struct {
	Port int
}

// ServerBuilder is responsible for building the server
type ServerBuilder struct {
	Config *ServerConfig
	Router router.Router
}

// NewServerBuilder creates a new ServerBuilder with default values
func NewServerBuilder() *ServerBuilder {
	port, err := strconv.Atoi(os.Getenv("PORT"))
	if err != nil {
		port = 8080
		log.Printf("error getting PORT from env. Err: %v. Set port to %v as default", err, port)
	}
	return &ServerBuilder{
		Config: &ServerConfig{Port: port},
		Router: nil,
	}
}

// WithRouter sets the router for the server explicitly
func (sb *ServerBuilder) WithRouter(r router.Router) *ServerBuilder {
	sb.Router = r
	return sb
}

// Build creates and returns the configured server
func (sb *ServerBuilder) Build() (*http.Server, error) {
	if sb.Router == nil {
		return nil, fmt.Errorf("router is required for the server")
	}
	handler := sb.Router.RegisterRoutes()

	// Declare Server config
	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", sb.Config.Port),
		Handler:      handler,
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}
	return server, nil
}
