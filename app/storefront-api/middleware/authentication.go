package middleware

import (
	"net/http"

	"github.com/anthonydip/flutter-messenger-go/app/storefront-api/webserver"

	"github.com/rs/zerolog/log"
)

// Authentication middleware
func Authentication(srv webserver.Server) func(h http.Handler) http.Handler {
	if srv == nil {
		log.Fatal().Msg("a nil dependency was passed to authentication middleware")
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Perform any authentication middleware

			next.ServeHTTP(w, r)
		})
	}
}
