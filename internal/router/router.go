package router

import (
	"net/http"

	bikes "bikesRentalAPI/internal/bikes/handlers"
	rentals "bikesRentalAPI/internal/rentals/handlers"
	users "bikesRentalAPI/internal/users/handlers"
	utils "bikesRentalAPI/internal/utils"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type Router interface {
	RegisterRoutes() http.Handler
}

type chiRouter struct {
	*chi.Mux
}

// New returns a new router interface using the chi router
func New() Router {
	r := chi.NewRouter()

	// Apply middleware
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.URLFormat)

	return &chiRouter{r}
}

// RegisterRoutes registers all routes for the application
func (r *chiRouter) RegisterRoutes() http.Handler {
	// Add routes here

	r.Get("/status", utils.StatusHandler) // TODO move to utils package

	r.Route("/users", func(r chi.Router) {
		// User authentication
		r.Post("/register", users.RegisterUser)
		r.Post("/login", users.LoginUser)
		r.Group(func(r chi.Router) {
			// User profile operations
			// r.Use(UserAuthMiddleware) // TODO implement Auth middleware
			r.Get("/profile", users.GetUserProfile)
			r.Patch("/profile", users.UpdateUserProfile)
		})
	})

	r.Route("/bikes", func(r chi.Router) {
		// Bike rental operations
		r.Group(func(r chi.Router) {
			// r.Use(UserAuthMiddleware)
			r.Get("/available", bikes.ListAvailableBikes)
			r.Post("/start", bikes.StartBikeRental)
			r.Post("/end", bikes.EndBikeRental)
			r.Get("/history", bikes.GetRentalHistory)
		})
	})

	r.Route("/admin", func(r chi.Router) {
		// Administrative endpoints
		// r.Use(AdminAuthMiddleware) // TODO implement Admin Auth middleware

		r.Route("/bikes", func(r chi.Router) {
			r.Post("/", bikes.AddBike)
			r.Patch("/{bike_id}", bikes.UpdateBike)
			r.Get("/", bikes.ListBikes)
		})

		r.Route("/users", func(r chi.Router) {
			r.Get("/", users.ListUsers)
			r.Get("/{user_id}", users.GetUserDetails)
			r.Patch("/{user_id}", users.UpdateUserDetails)
		})

		r.Route("/rentals", func(r chi.Router) {
			r.Get("/", rentals.ListRentals)
			r.Get("/{rental_id}", rentals.GetRentalDetails)
			r.Patch("/{rental_id}", rentals.UpdateRentalDetails)
		})
	})
	return r
}
