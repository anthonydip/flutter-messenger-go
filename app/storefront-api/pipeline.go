package main

import (
	"net/http"

	"github.com/anthonydip/flutter-messenger-go/app/storefront-api/middleware"
	"github.com/anthonydip/flutter-messenger-go/app/storefront-api/routes"
	"github.com/anthonydip/flutter-messenger-go/app/storefront-api/routes/auth/signin"
	accessToken "github.com/anthonydip/flutter-messenger-go/app/storefront-api/routes/auth/tokens/access"
	"github.com/anthonydip/flutter-messenger-go/app/storefront-api/routes/users"
	"github.com/anthonydip/flutter-messenger-go/app/storefront-api/webserver"
	"github.com/gorilla/mux"
	"github.com/rs/zerolog/log"
)

// Build the HTTP pipeline
func BuildPipeline(srv webserver.Server, r *mux.Router) {
	log.Info().Msg("building pipeline...")
	r.HandleFunc("/ping", routes.Ping(srv)).Methods(http.MethodGet)

	r.Use(middleware.Authentication(srv))

	r.HandleFunc("/auth/signin", signin.Post(srv)).Methods(http.MethodPost)

	r.HandleFunc("/auth/tokens/access", accessToken.Post(srv)).Methods(http.MethodPost)

	r.HandleFunc("/users/{userID}", users.Get(srv)).Methods(http.MethodGet)
	r.HandleFunc("/users", users.Post(srv)).Methods(http.MethodPost)
}
