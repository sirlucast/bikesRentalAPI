package server

import (
	"net/http"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	portMock = "1111"
)

// MockRouter implements the router.Router interface for testing purposes
type MockRouter struct{}

func (mr *MockRouter) RegisterRoutes() http.Handler {
	// Mock implementation for testing
	return http.NewServeMux()
}

func TestServerBuilder(t *testing.T) {
	// GIVEN: a mock port
	t.Setenv("PORT", portMock)
	t.Run("Success - NewServerBuilder is built successfully", func(t *testing.T) {
		// GIVEN: a server builder
		builder := NewServerBuilder()
		// WHEN: the server builder is created
		// THEN: the server builder should be created successfully
		assert.NotNil(t, builder)
		// THEN: the server builder should have the default port got from Env.
		assert.Equal(t, portMock, strconv.Itoa(builder.Config.Port))
	})
	t.Run("Success - WithRouter sets the router for the server", func(t *testing.T) {
		// GIVEN: a server builder and a mock router
		builder := NewServerBuilder()

		handler := http.NewServeMux()
		// WHEN: the router is set for the server
		builder.WithHanlder(handler)
		// THEN: the router should be set for the server
		assert.Equal(t, handler, builder.Handler)
	})
	t.Run("Success - Build creates and returns the configured server", func(t *testing.T) {
		// GIVEN: a server builder and a mock router
		builder := NewServerBuilder()
		handler := http.NewServeMux()
		builder.WithHanlder(handler)
		// WHEN: the server is built
		server, err := builder.Build()
		// THEN: the server should be built successfully
		assert.NoError(t, err)
		assert.NotNil(t, server)
	})
}
