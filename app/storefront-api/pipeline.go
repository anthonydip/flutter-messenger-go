package main

import (
	"net/http"

	"github.com/anthonydip/flutter-messenger-go/app/storefront-api/middleware"
	"github.com/anthonydip/flutter-messenger-go/app/storefront-api/routes/users"
	"github.com/anthonydip/flutter-messenger-go/app/storefront-api/webserver"
	"github.com/gorilla/mux"
	"github.com/rs/zerolog/log"
)

// Build the HTTP pipeline
func BuildPipeline(srv webserver.Server, r *mux.Router) {
	log.Info().Msg("building pipeline...")

	r.Use(middleware.Authentication(srv))

	r.HandleFunc("/users/{userID}", users.Get(srv)).Methods(http.MethodGet)
}
