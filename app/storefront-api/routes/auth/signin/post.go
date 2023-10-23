package signin

import (
	"encoding/json"
	"net/http"

	"github.com/anthonydip/flutter-messenger-go/app/storefront-api/utils"
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

// Sign a user in using email and password
func Post(srv webserver.Server) http.HandlerFunc {
	if srv == nil {
		log.Fatal().Msg("a nil dependency was passed to POST '/auth/signin'")
	}

	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		user := dtos.User{}

		r.Body = http.MaxBytesReader(w, r.Body, 1048576)

		dec := json.NewDecoder(r.Body)
		dec.DisallowUnknownFields()

		err := dec.Decode(&user)
		if err != nil {
			log.Error().Msg("[POST /auth/signin] Unable to decode user")

			w.WriteHeader(http.StatusBadRequest)

			res := Response{
				Status:        "BAD REQUEST",
				StatusCode:    400,
				StatusMessage: "Invalid request body",
			}

			json.NewEncoder(w).Encode(&res)
			return
		}

		log.Info().Msgf("[POST /auth/signin] Received a request, %+v", user)

		sublogger := log.With().Any("request", user).Logger()

		// Validate user sign in request
		err = utils.ValidatePostUser(user)
		if err != nil {
			res := Response{
				Status:     "BAD REQUEST",
				StatusCode: 400,
			}
			w.WriteHeader(http.StatusBadRequest)

			switch err.Error() {
			case "invalid email":
				sublogger.Error().Msgf("[POST /auth/signin] Invalid email, received %s", user.Email)
				res.StatusMessage = "Invalid email, criteria not met"
			case "invalid provider":
				sublogger.Error().Msgf("[POST /auth/signin] Invalid provider, received %s", user.Provider)
				res.StatusMessage = "Invalid provider, criteria not met"
			case "invalid password":
				sublogger.Error().Msg("[POST /auth/signin] Invalid password")
				res.StatusMessage = "Invalid password, criteria not met"
			default:
				sublogger.Error().Msgf("[POST /auth/signin] Error occurred validating user, %v", err)
				res = Response{
					Status:        "INTERNAL SERVER ERROR",
					StatusCode:    500,
					StatusMessage: "Error validating user request",
				}
				w.WriteHeader(http.StatusInternalServerError)
			}

			json.NewEncoder(w).Encode(&res)
			return
		}

		// Attempt to sign the user in by validating password with hash
		err = srv.SignIn(user)
		if err != nil {
			res := Response{
				Status:     "UNAUTHORIZED",
				StatusCode: 401,
			}
			w.WriteHeader(http.StatusBadRequest)

			switch err.Error() {
			case "invalid password":
				sublogger.Error().Msg("[POST /auth/signin] Password does not match hash")
				res.StatusMessage = "Incorrect password"
			default:
				sublogger.Error().Msgf("[POST /auth/signin] Error occurred signing user in, %v", err)
				res = Response{
					Status:        "INTERNAL SERVER ERROR",
					StatusCode:    500,
					StatusMessage: "Error validating user request",
				}
				w.WriteHeader(http.StatusInternalServerError)
			}

			json.NewEncoder(w).Encode(&res)
			return
		}

		// Generate access token for the user
		token, err := srv.GenerateAccessToken(user)
		if err != nil {
			switch err.Error() {
			case "error reading file":
				sublogger.Error().Msg("[POST /auth/signin] Error reading private key")
			case "error parsing pem":
				sublogger.Error().Msg("[POST /auth/signin] Error parsing PEM file")
			case "error signing token":
				sublogger.Error().Msg("[POST /auth/signin] Error signing token")
			default:
				sublogger.Error().Msg("[POST /auth/signin] Unexpected error occurred generating user access token")
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

		sublogger.Info().Msgf("[POST /auth/signin] Successfully signed user in with access token %s", token)

		w.WriteHeader(http.StatusOK)

		res := Response{
			Status:        "SUCCESS",
			StatusCode:    200,
			StatusMessage: "Successfully signed user in",
			Token:         token,
		}

		json.NewEncoder(w).Encode(&res)
	}
}
