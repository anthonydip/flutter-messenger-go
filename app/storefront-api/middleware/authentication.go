package middleware

import (
	"encoding/json"
	"net/http"
	"regexp"
	"strings"

	"github.com/anthonydip/flutter-messenger-go/app/storefront-api/webserver"

	"github.com/rs/zerolog/log"
)

type Response struct {
	Status        string `json:"status"`
	StatusCode    int    `json:"statusCode"`
	StatusMessage string `json:"statusMessage,omitempty"`
}

type route struct {
	Regex  string
	Method string
}

var internalRoutes = [...]route{
	{Regex: "^/auth/signin$", Method: http.MethodPost},
	{Regex: "^/auth/tokens/access$", Method: http.MethodPost},
	{Regex: "^/auth/tokens/access/", Method: http.MethodDelete},
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
			w.Header().Set("Content-Type", "application/json")

			token := r.Header.Get("Authorization")
			tokenString := strings.ReplaceAll(token, "Bearer ", "")

			// Check if the request URL requires an internal token
			isInternalRoute := false
			for _, route := range internalRoutes {
				match, _ := regexp.MatchString(route.Regex, r.URL.String())
				if match {
					isInternalRoute = true
				}
			}

			res := Response{
				Status:        "UNAUTHORIZED",
				StatusCode:    401,
				StatusMessage: "Invalid authorization token",
			}

			// Perform authentication middleware depending on the route
			if isInternalRoute {
				if !srv.ValidateInternalJWT(tokenString) {
					w.WriteHeader(http.StatusUnauthorized)
					json.NewEncoder(w).Encode(&res)
					return
				}
			} else {
				if !srv.ValidateJWT(tokenString) {
					w.WriteHeader(http.StatusUnauthorized)
					json.NewEncoder(w).Encode(&res)
					return
				}
			}

			next.ServeHTTP(w, r)
		})
	}
}
