package router

import (
	bikemocks "bikesRentalAPI/internal/bikes/handlers/mocks"
	rentalmocks "bikesRentalAPI/internal/rentals/handlers/mocks"
	usermocks "bikesRentalAPI/internal/users/handlers/mocks"

	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

const (
	statusURL = "/status"
)

func TestRouter(t *testing.T) {
	t.Setenv("ADMIN_CREDENTIALS", "dGVzdDp0ZXN0") // test:test
	// GIVEN a mock user handler
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockUserHandler := usermocks.NewMockHandler(mockCtrl)
	mockBikeHandler := bikemocks.NewMockHandler(mockCtrl)
	mockRentalHandler := rentalmocks.NewMockHandler(mockCtrl)

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
		handler := router.RegisterRoutes(mockUserHandler, mockBikeHandler, mockRentalHandler)
		// THEN: the handler should be created successfully
		assert.NotNil(t, handler)
	})
	t.Run("Success - statusHandler returns a status message: '.'", func(t *testing.T) {
		// GIVEN: a router, a test server and expected message
		router := New()
		// WHEN: the routes are registered
		handler := router.RegisterRoutes(mockUserHandler, mockBikeHandler, mockRentalHandler)
		// GIVEN: a test server
		server := httptest.NewServer(handler)
		defer server.Close()
		expectedMessage := "." // default status message from Heartbeats middleware
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
