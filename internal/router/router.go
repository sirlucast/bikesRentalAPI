package router

import (
	"encoding/json"
	"log"
	"net/http"

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
	r.Route("/", func(r chi.Router) {
		r.Get("/status", statusHandler)
	})
	return r
}

func statusHandler(w http.ResponseWriter, r *http.Request) {
	resp := map[string]string{"message": "Server is up and running"}
	jsonResp, err := json.Marshal(resp)
	if err != nil {
		log.Printf("error marshalling response: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(jsonResp)
}
