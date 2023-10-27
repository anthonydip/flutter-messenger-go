package friends

import (
	"encoding/json"
	"net/http"

	"github.com/anthonydip/flutter-messenger-go/app/storefront-api/webserver"

	"github.com/anthonydip/flutter-messenger-go/app/storefront-api/utils"
	"github.com/rs/zerolog/log"
)

type FriendRequest struct {
	Email string `json:"email,omitempty"`
}

type FriendResponse struct {
	Status        string `json:"status"`
	StatusCode    int    `json:"statusCode"`
	StatusMessage string `json:"statusMessage,omitempty"`
}

// Add a friend
func Post(srv webserver.Server) http.HandlerFunc {
	if srv == nil {
		log.Fatal().Msg("a nil dependency was passed to POST '/users/friends'")
	}

	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		friend := FriendRequest{}

		r.Body = http.MaxBytesReader(w, r.Body, 1048576)

		dec := json.NewDecoder(r.Body)
		dec.DisallowUnknownFields()

		err := dec.Decode(&friend)
		if err != nil {
			log.Error().Msg("[POST /users/friends] Unable to decode friend")

			w.WriteHeader(http.StatusBadRequest)

			res := FriendResponse{
				Status:        "BAD REQUEST",
				StatusCode:    400,
				StatusMessage: "Invalid request body",
			}

			json.NewEncoder(w).Encode(&res)
			return
		}

		log.Info().Msgf("[POST /users/friends] Received a request: %+v", friend)

		sublogger := log.With().Any("request", friend).Logger()

		// Validate the email
		err = utils.ValidateEmail(friend.Email)
		if err != nil {
			sublogger.Error().Msgf("[POST /users/friends] Invalid email provided")
			res := FriendResponse{
				Status:        "BAD REQUEST",
				StatusCode:    400,
				StatusMessage: "Invalid email",
			}

			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(&res)
			return
		}

		// Extract the email from the access token
		authHeader := r.Header.Get("Authorization")
		token, err := utils.GetAuthorizationToken(authHeader)
		if err != nil {
			res := FriendResponse{
				Status:     "UNAUTHORIZED",
				StatusCode: 401,
			}

			switch err.Error() {
			case "empty header":
				sublogger.Error().Msgf("[POST /users/friends] Empty authorization header provided")
				res.StatusMessage = "Empty authorization header"
			case "invalid header":
				sublogger.Error().Msgf("[POST /users/friends] Invalid authorization header provided")
				res.StatusMessage = "Invalid authorization header"
			default:
				res = FriendResponse{
					Status:        "INTERNAL SERVER ERROR",
					StatusCode:    500,
					StatusMessage: "Error occurred extracting authorization token",
				}
				w.WriteHeader(http.StatusInternalServerError)
				json.NewEncoder(w).Encode(&res)
				return
			}

			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(&res)
			return
		}

		// Parse the user from the authorization token
		user, err := srv.ValidateParseJWT(token)
		if err != nil {
			switch err.Error() {
			case "error reading pem":
				sublogger.Error().Msgf("[POST /users/friends] Error reading PEM for token")
			case "error parsing pem":
				sublogger.Error().Msgf("[POST /users/friends] Error parsing PEM for token")
			case "invalid token":
				res := FriendResponse{
					Status:        "UNAUTHORIZED",
					StatusCode:    401,
					StatusMessage: "Invalid authorization token",
				}
				w.WriteHeader(http.StatusUnauthorized)
				json.NewEncoder(w).Encode(&res)
				return
			default:
				sublogger.Error().Msgf("[POST /users/friends] Error occurred validating and parsing token")
			}

			res := FriendResponse{
				Status:        "INTERNAL SERVER ERROR",
				StatusCode:    500,
				StatusMessage: "Error validating and parsing token",
			}
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(&res)
			return
		}

		// Retrieve the friend from the email
		friendUser, err := srv.GetUserByEmail(friend.Email)
		if err != nil {
			if err.Error() == "user does not exist" {
				sublogger.Info().Msgf("[POST /users/friends] Friend does not exist")
				res := FriendResponse{
					Status:        "NOT FOUND",
					StatusCode:    404,
					StatusMessage: "Friend does not exist",
				}
				w.WriteHeader(http.StatusNotFound)
				json.NewEncoder(w).Encode(&res)
				return
			} else {
				sublogger.Error().Msgf("[POST /users/friends] Error retrieving friend from the database")
				res := FriendResponse{
					Status:        "INTERNAL SERVER ERROR",
					StatusCode:    500,
					StatusMessage: "Error retrieving friend",
				}
				w.WriteHeader(http.StatusInternalServerError)
				json.NewEncoder(w).Encode(&res)
				return
			}
		}

		// Attempt to add the friend
		err = srv.PostFriend(user.Id, friendUser)
		if err != nil {
			sublogger.Error().Msgf("[POST /users/friends] Error adding friend, %s", err.Error())
			res := FriendResponse{
				Status:        "INTERNAL SERVER ERROR",
				StatusCode:    500,
				StatusMessage: "Error adding friend",
			}
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(&res)
			return
		}

		sublogger.Info().Msgf("[POST /users/friends] Successfully added friend")

		res := FriendResponse{
			Status:        "SUCCESS",
			StatusCode:    200,
			StatusMessage: "Successfully added friend",
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(&res)
	}
}
