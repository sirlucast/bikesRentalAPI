package router

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	statusURL = "/status"
)

func TestRouter(t *testing.T) {
	t.Run("Success - New returns a new router interface using the chi router", func(t *testing.T) {
		// WHEN: the New function is called
		router := New()
		// THEN: the router should be created successfully
		assert.NotNil(t, router)
	})
	t.Run("Success - RegisterRoutes registers all routes for the application", func(t *testing.T) {
		// GIVEN: a router
		router := New()
		// WHEN: the routes are registered
		handler := router.RegisterRoutes()
		// THEN: the handler should be created successfully
		assert.NotNil(t, handler)
	})
	t.Run("Success - statusHandler returns a status message", func(t *testing.T) {
		// GIVEN: a router, a test server and expected message
		router := New()
		server := httptest.NewServer(router.RegisterRoutes())
		defer server.Close()
		expectedMessage := "Server is up and running"
		// WHEN: a request is made to the server
		resp, err := http.Get(server.URL + statusURL)
		assert.NoError(t, err)
		defer resp.Body.Close()
		// THEN: the response should be OK
		assert.NotNil(t, resp)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
		// THEN: the response body should contain the status message
		body, err := io.ReadAll(resp.Body)
		assert.NoError(t, err)
		assert.Contains(t, string(body), expectedMessage)

	})
}
