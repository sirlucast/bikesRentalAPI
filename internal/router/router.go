package router

import (
	"log"
	"net/http"
	"os"
	"strings"

	bikes "bikesRentalAPI/internal/bikes/handlers"
	"bikesRentalAPI/internal/helpers"
	rentals "bikesRentalAPI/internal/rentals/handlers"
	users "bikesRentalAPI/internal/users/handlers"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/jwtauth/v5"
)

type Router interface {
	RegisterRoutes(userHandler users.Handler) http.Handler
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
	r.Use(middleware.Heartbeat("/status"))

	return &chiRouter{r}
}

// RegisterRoutes registers all routes for the application
func (r *chiRouter) RegisterRoutes(userHandler users.Handler) http.Handler {
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
			r.Get("/profile", users.GetUserProfile)
			r.Patch("/profile", users.UpdateUserProfile)
		})
	})

	r.Route("/bikes", func(r chi.Router) {
		// Bike rental operations
		r.Group(func(r chi.Router) {
			r.Use(jwtauth.Verifier(tokenAuth))
			r.Use(jwtauth.Authenticator(tokenAuth))
			r.Get("/available", bikes.ListAvailableBikes)
			r.Post("/start", bikes.StartBikeRental)
			r.Post("/end", bikes.EndBikeRental)
			r.Get("/history", bikes.GetRentalHistory)
		})
	})

	r.Route("/admin", func(r chi.Router) {
		// Administrative endpoints
		r.Group(func(r chi.Router) {
			r.Use(middleware.BasicAuth("bikesRental API Administration", adminCredentials))
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
	})
	return r
}

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
