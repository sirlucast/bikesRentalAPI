package router

import (
	"log"
	"net/http"
	"os"
	"strings"

	bikes "bikesRentalAPI/internal/bikes/handlers"
	"bikesRentalAPI/internal/helpers"
	"bikesRentalAPI/internal/middlewares"
	rentals "bikesRentalAPI/internal/rentals/handlers"
	users "bikesRentalAPI/internal/users/handlers"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/jwtauth/v5"
)

type Router interface {
	RegisterRoutes(userHandler users.Handler, bikeHandler bikes.Handler, rentalHandler rentals.Handler) http.Handler
}

type chiRouter struct {
	*chi.Mux
}

var (
	tokenAuth           *jwtauth.JWTAuth
	adminCredentials    map[string]string
	secretJWTKey        string
	adminCredentialsEnc string
)

func init() {
	secretJWTKey = os.Getenv("JWT_SECRET_KEY")
	adminCredentialsEnc = os.Getenv("ADMIN_CREDENTIALS")
	tokenAuth = jwtauth.New("HS256", []byte(secretJWTKey), nil)
	adminCredentials = adminCredentialsDecode()
}

// New returns a new router interface using the chi router
func New() Router {
	r := chi.NewRouter()

	// Apply middleware
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.URLFormat)
	r.Use(middleware.RequestID)
	r.Use(middleware.Heartbeat("/status"))

	return &chiRouter{r}
}

// RegisterRoutes registers all routes for the application
func (r *chiRouter) RegisterRoutes(userHandler users.Handler, bikeHandler bikes.Handler, rentalHandler rentals.Handler) http.Handler {
	// Add routes here
	r.Route("/users", func(r chi.Router) {
		// User authentication
		r.Post("/register", userHandler.RegisterUser)
		r.Post("/login", func(w http.ResponseWriter, r *http.Request) {
			userHandler.LoginUser(tokenAuth, w, r)
		})
		r.Group(func(r chi.Router) {
			r.Use(jwtauth.Verifier(tokenAuth))
			r.Use(jwtauth.Authenticator(tokenAuth))
			// User profile operations
			r.Get("/profile", userHandler.GetUserProfile)
			r.Patch("/profile", userHandler.UpdateUserProfile)
		})
	})

	r.Route("/bikes", func(r chi.Router) {
		r.Group(func(r chi.Router) {
			r.Use(jwtauth.Verifier(tokenAuth))
			r.Use(jwtauth.Authenticator(tokenAuth))
			// Bike general operations
			r.With(middlewares.Pagination).Get("/available", bikeHandler.ListAvailableBikes)
		})
	})
	r.Route("/rentals", func(r chi.Router) {
		r.Group(func(r chi.Router) {
			r.Use(jwtauth.Verifier(tokenAuth))
			r.Use(jwtauth.Authenticator(tokenAuth))
			// Rental operations
			r.Post("/start", rentalHandler.StartBikeRental)
			r.Post("/end", rentalHandler.EndBikeRental)
			r.With(middlewares.Pagination).Get("/history", rentalHandler.GetRentalHistoryByUserID)
		})
	})

	r.Route("/admin", func(r chi.Router) {
		// Administrative endpoints
		r.Group(func(r chi.Router) {
			r.Use(middleware.BasicAuth("bikesRental API Administration", adminCredentials))

			r.Route("/bikes", func(r chi.Router) {
				r.Post("/", bikeHandler.AddBike)
				r.Patch("/{bike_id}", bikeHandler.UpdateBike)
				r.Get("/{bike_id}", bikeHandler.GetBikeByID)
				r.With(middlewares.Pagination).Get("/", bikeHandler.ListAllBikes)
			})

			r.Route("/users", func(r chi.Router) {
				r.With(middlewares.Pagination).Get("/", userHandler.ListAllUsers)
				r.Get("/{user_id}", userHandler.GetUserDetails)
				r.Patch("/{user_id}", userHandler.UpdateUserDetails)
			})

			r.Route("/rentals", func(r chi.Router) {
				r.With(middlewares.Pagination).Get("/", rentalHandler.GetRentalList)
				r.Get("/{rental_id}", rentalHandler.GetRentalDetails)
				r.Patch("/{rental_id}", rentalHandler.UpdateRentalDetails)
			})
		})
	})
	return r
}

// adminCredentialsDecode decodes the admin credentials from the environment variable
func adminCredentialsDecode() map[string]string {
	decodedAdminCred, err := helpers.Base64Decode(adminCredentialsEnc)
	if err {
		log.Printf("failed to decode admin credentials.")
		return nil
	}
	credentials := strings.Split(decodedAdminCred, ":")
	if len(credentials) != 2 {
		log.Printf("decoded admin credentials are not following <user:passowrd> shape. Got: %v", credentials)
		return nil
	}
	return map[string]string{credentials[0]: credentials[1]}
}
