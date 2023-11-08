package main

import (
	"net/http"

	"github.com/anthonydip/flutter-messenger-go/app/storefront-api/middleware"
	"github.com/anthonydip/flutter-messenger-go/app/storefront-api/routes"
	"github.com/anthonydip/flutter-messenger-go/app/storefront-api/routes/auth/signin"
	accessToken "github.com/anthonydip/flutter-messenger-go/app/storefront-api/routes/auth/tokens/access"
	"github.com/anthonydip/flutter-messenger-go/app/storefront-api/routes/users"
	"github.com/anthonydip/flutter-messenger-go/app/storefront-api/routes/users/friends"
	"github.com/anthonydip/flutter-messenger-go/app/storefront-api/webserver"
	"github.com/anthonydip/flutter-messenger-go/app/storefront-api/ws"
	"github.com/gorilla/mux"
	"github.com/rs/zerolog/log"
)

// Build the HTTP pipeline
func BuildPipeline(srv webserver.Server, hub *ws.Hub, r *mux.Router) {
	log.Info().Msg("building pipeline...")
	r.HandleFunc("/ping", routes.Ping(srv)).Methods(http.MethodGet)

	r.Use(middleware.Authentication(srv))

	r.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		log.Info().Msgf("[GET /ws] Established WebSocket connection for %s", r.URL.Query().Get("token"))
		ws.ServeWs(hub, w, r)
	})

	r.HandleFunc("/auth/signin", signin.Post(srv)).Methods(http.MethodPost)

	r.HandleFunc("/auth/tokens/access", accessToken.Post(srv)).Methods(http.MethodPost)
	r.HandleFunc("/auth/tokens/access/{token}", accessToken.Delete(srv)).Methods(http.MethodDelete)

	r.HandleFunc("/users/{id:(?:[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}|[^@]+@[^/]+)}", users.Get(srv)).Methods(http.MethodGet)
	r.HandleFunc("/users", users.Post(srv)).Methods(http.MethodPost)
	r.HandleFunc("/users/friends", friends.Post(srv)).Methods(http.MethodPost)
	r.HandleFunc("/users/friends", friends.Get(srv)).Methods(http.MethodGet)
}
