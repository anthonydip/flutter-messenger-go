package access

import (
	"encoding/json"
	"net/http"

	"github.com/anthonydip/flutter-messenger-go/app/storefront-api/webserver"
	"github.com/anthonydip/flutter-messenger-go/pkg/dtos"
	"github.com/rs/zerolog/log"
)

type Response struct {
	Status        string `json:"status"`
	StatusCode    int    `json:"statusCode"`
	StatusMessage string `json:"statusMessage,omitempty"`
	Token         string `json:"token,omitempty"`
}

// Retrieve a user access token
func Post(srv webserver.Server) http.HandlerFunc {
	if srv == nil {
		log.Fatal().Msg("a nil dependency was passed to GET '/tokens/access'")
	}

	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		user := dtos.User{}

		r.Body = http.MaxBytesReader(w, r.Body, 1048576)

		dec := json.NewDecoder(r.Body)
		dec.DisallowUnknownFields()

		err := dec.Decode(&user)
		if err != nil {
			log.Error().Msg("[POST /tokens/access] Unable to decode user")

			w.WriteHeader(http.StatusBadRequest)

			res := Response{
				Status:        "BAD REQUEST",
				StatusCode:    400,
				StatusMessage: "Invalid request body",
			}

			json.NewEncoder(w).Encode(&res)
			return
		}

		log.Info().Msgf("[POST /tokens/access] Received a request, %+v", user)

		// Check if the user exists
		user, err = srv.GetUserByEmail(user.Email)
		if err != nil {
			log.Info().Msgf("[POST /tokens/access] Specified user does not exist")

			res := Response{
				Status:        "NOT FOUND",
				StatusCode:    404,
				StatusMessage: "User does not exist",
			}

			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(&res)
			return
		}

		// Generate user access token
		token, err := srv.GenerateAccessToken(user)
		if err != nil {
			switch err.Error() {
			case "error reading file":
				log.Error().Msg("[POST /tokens/access] Error reading private key")
			case "error parsing pem":
				log.Error().Msg("[POST /tokens/access] Error parsing PEM file")
			case "error signing token":
				log.Error().Msg("[POST /tokens/access] Error signing token")
			default:
				log.Error().Msg("[POST /tokens/access] Unexpected error occurred generating user access token")
			}

			res := Response{
				Status:        "INTERNAL SERVER ERROR",
				StatusCode:    500,
				StatusMessage: "Error generating access token",
			}

			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(&res)
			return
		}

		res := Response{
			Status:        "CREATED",
			StatusCode:    201,
			StatusMessage: "Successfully generated access token",
			Token:         token,
		}

		log.Info().Msgf("[POST /tokens/access] Successfully generated user access token, %s", token)

		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(&res)
	}
}
