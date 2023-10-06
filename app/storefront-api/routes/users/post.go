package users

import (
	"encoding/json"
	"net/http"

	"github.com/anthonydip/flutter-messenger-go/app/storefront-api/utils"
	"github.com/anthonydip/flutter-messenger-go/app/storefront-api/webserver"
	"github.com/anthonydip/flutter-messenger-go/pkg/dtos"

	"github.com/rs/zerolog/log"
)

type Response struct {
	Status        string     `json:"status"`
	StatusCode    int        `json:"statusCode"`
	StatusMessage string     `json:"statusMessage,omitempty"`
	User          *dtos.User `json:"user,omitempty"`
}

// Create a user
func Post(srv webserver.Server) http.HandlerFunc {
	if srv == nil {
		log.Fatal().Msg("A nil dependency was passed to POST '/users'")
	}

	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		user := dtos.User{}

		r.Body = http.MaxBytesReader(w, r.Body, 1048576)

		dec := json.NewDecoder(r.Body)
		dec.DisallowUnknownFields()

		err := dec.Decode(&user)
		if err != nil {
			log.Error().Msg("[POST /users] Unable to decode user")

			w.WriteHeader(http.StatusBadRequest)

			res := Response{
				Status:        "BAD REQUEST",
				StatusCode:    400,
				StatusMessage: "Invalid request body",
			}

			json.NewEncoder(w).Encode(&res)
			return
		}

		log.Info().Msgf("[POST /users] Received a request, %+v", user)

		// Create sub-logger
		sublogger := log.With().Any("request", user).Logger()

		if err != nil {
			sublogger.Error().Msg("[POST /users] Invalid request body")

			w.WriteHeader(http.StatusBadRequest)

			res := Response{
				Status:        "BAD REQUEST",
				StatusCode:    400,
				StatusMessage: "Invalid request body",
			}

			json.NewEncoder(w).Encode(&res)
			return
		}

		// Validate the user request
		err = utils.ValidatePostUser(user)
		if err != nil {
			res := Response{
				Status:     "BAD REQUEST",
				StatusCode: 400,
			}
			w.WriteHeader(http.StatusBadRequest)

			switch err.Error() {
			case "invalid email":
				sublogger.Error().Msgf("[POST /users] Invalid email, received %s", user.Email)
				res.StatusMessage = "Invalid email"
			case "invalid provider":
				sublogger.Error().Msgf("[POST /users] Invalid provider, received %s", user.Provider)
				res.StatusMessage = "Invalid provider"
			case "invalid password":
				sublogger.Error().Msg("[POST /users] Invalid password")
				res.StatusMessage = "Invalid password, criteria not met"
			default:
				sublogger.Error().Msgf("[POST /users] Error occurred validating user request, %v", err)
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

		result, err := srv.PostUser(user)

		if err != nil {
			res := Response{}

			if err.Error() == "409 Conflict" {
				sublogger.Error().Msgf("[POST /users] User already exists with email: %s", user.Email)

				res = Response{
					Status:        "CONFLICT",
					StatusCode:    409,
					StatusMessage: "User already exists",
				}

				w.WriteHeader(http.StatusConflict)
				json.NewEncoder(w).Encode(&res)
			} else {
				sublogger.Error().Msgf("[POST /users] Posting new user failed, %v", err.Error())

				res = Response{
					Status:     "INTERNAL SERVER ERROR",
					StatusCode: 500,
				}

				w.WriteHeader(http.StatusInternalServerError)
				json.NewEncoder(w).Encode(&res)
			}
			return
		}

		sublogger.Info().Msgf("[POST /users] Successfully created user: %+v", result)

		res := Response{
			Status:        "CREATED",
			StatusCode:    201,
			StatusMessage: "User successfully created",
			User:          &result,
		}

		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(&res)
	}
}
