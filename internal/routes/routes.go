package routes

import (
	"net/http"

	"github.com/edaywalid/url-shortner/internal/handlers"
	"github.com/gorilla/mux"
	"github.com/rs/zerolog/log"
)

type Routes struct {
	handlers *handlers.Handler
}

func NewRoutes(handlers *handlers.Handler) *Routes {
	return &Routes{handlers}
}

func (r *Routes) RegisterRoutes() {
	log.Info().Msg("Registering routes")
	router := mux.NewRouter()
	router.HandleFunc("/shorten", r.handlers.GetShortURL).Methods("POST")
	router.HandleFunc("/{shortCode}", r.handlers.Redirect).Methods("GET")

	log.Info().Msg("Routes registered")

	log.Info().Msg("Starting the server on port 8080")
	if err := http.ListenAndServe(":8080", router); err != nil {
		log.Fatal().Err(err).Msg("Failed to start the server")
	}
}
