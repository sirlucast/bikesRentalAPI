package server

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	_ "github.com/joho/godotenv/autoload"
)

// ServerConfig holds configuration parameters for the server
type ServerConfig struct {
	Port int
}

// ServerBuilder is responsible for building the server
type ServerBuilder struct {
	Config  *ServerConfig
	Handler http.Handler
}

// NewServerBuilder creates a new ServerBuilder with default values
func NewServerBuilder() *ServerBuilder {
	port, err := strconv.Atoi(os.Getenv("PORT"))
	if err != nil {
		port = 8080
		log.Printf("error getting PORT from env. Err: %v. Set port to %v as default", err, port)
	}
	return &ServerBuilder{
		Config:  &ServerConfig{Port: port},
		Handler: nil,
	}
}

// WithHanlder sets the handler for the server explicitly
func (sb *ServerBuilder) WithHanlder(r http.Handler) *ServerBuilder {
	sb.Handler = r
	return sb
}

// Build creates and returns the configured server
func (sb *ServerBuilder) Build() (*http.Server, error) {
	if sb.Handler == nil {
		return nil, fmt.Errorf("router is required for the server")
	}
	//handler := sb.Router.RegisterRoutes()

	// Declare Server config
	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", sb.Config.Port),
		Handler:      sb.Handler,
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}
	return server, nil
}
