package middleware

import (
	"net/http"
	"regexp"

	"github.com/anthonydip/flutter-messenger-go/app/storefront-api/webserver"

	"github.com/rs/zerolog/log"
)

type route struct {
	Regex  string
	Method string
}

var internalRoutes = [...]route{
	{Regex: "^/tokens/access$", Method: http.MethodGet},
	{Regex: "^/users$", Method: http.MethodPost},
	{Regex: "^/users/", Method: http.MethodGet},
}

// Authentication middleware
func Authentication(srv webserver.Server) func(h http.Handler) http.Handler {
	if srv == nil {
		log.Fatal().Msg("a nil dependency was passed to authentication middleware")
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			token := r.Header.Get("Authorization")

			isInternalRoute := false
			// Check if the request URL requires an internal token
			for _, route := range internalRoutes {
				match, _ := regexp.MatchString(route.Regex, r.URL.String())
				if match {
					isInternalRoute = true
				}
			}

			// Perform authentication middleware depending on the route
			if isInternalRoute {
				if !srv.ValidateInternalJWT(token) {
					w.WriteHeader(http.StatusUnauthorized)
					w.Write([]byte("401 Unauthorized"))
					return
				}
			} else {
				if !srv.ValidateJWT(token) {
					w.WriteHeader(http.StatusUnauthorized)
					w.Write([]byte("401 Unauthorized"))
					return
				}
			}

			next.ServeHTTP(w, r)
		})
	}
}
