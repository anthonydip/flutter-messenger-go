package users

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/anthonydip/flutter-messenger-go/app/storefront-api/webserver"

	"github.com/gorilla/mux"
	"github.com/rs/zerolog/log"
)

// Get a user
func Get(srv webserver.Server) http.HandlerFunc {
	if srv == nil {
		log.Fatal().Msg("a nil dependency was passed to GET '/users/{userID}'")
	}

	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		// Get the userID
		params := mux.Vars(r)
		userID := strings.TrimSpace(params["userID"])

		sublogger := log.With().Any("userID", userID).Logger()
		sublogger.Info().Msg("[GET /users/{userID}] Received a request")

		// Get the query type
		idType := r.URL.Query().Get("type")

		// Default idType to "id" if empty or does not exist
		if idType == "" {
			idType = "id"
		}

		if idType != "email" && idType != "id" {
			sublogger.Info().Msgf("[GET /users/{userID}] Invalid idType %s", idType)
			res := Response{
				Status:        "BAD REQUEST",
				StatusCode:    400,
				StatusMessage: "Invalid ID type in query parameter",
			}
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(&res)
			return
		}

		if idType == "id" {
			// Get the user from the userID
			user, err := srv.GetUser(userID)
			if err != nil {
				sublogger.Info().Msg("[GET /users/{userID}] User does not exist")

				res := Response{
					Status:        "NOT FOUND",
					StatusCode:    404,
					StatusMessage: "User does not exist",
				}
				w.WriteHeader(http.StatusNotFound)
				json.NewEncoder(w).Encode(&res)
				return
			}

			sublogger.Info().Msgf("[GET /users/{userID}] Successfully retrieved user: %+v", user)

			res := Response{
				Status:        "SUCCESS",
				StatusCode:    200,
				StatusMessage: "User exists",
				User:          &user,
			}
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(&res)
		} else {
			// Get the user from email
			user, err := srv.GetUserByEmail(userID)
			if err != nil {
				sublogger.Info().Msg("[GET /users/{userID}] User does not exist")

				res := Response{
					Status:        "NOT FOUND",
					StatusCode:    404,
					StatusMessage: "User does not exist",
				}
				w.WriteHeader(http.StatusNotFound)
				json.NewEncoder(w).Encode(&res)
				return
			}

			sublogger.Info().Msgf("[GET /users/{userID}] Successfully retrieved user: %+v", user)

			res := Response{
				Status:        "SUCCESS",
				StatusCode:    200,
				StatusMessage: "User exists",
				User:          &user,
			}
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(&res)
		}

	}
}
