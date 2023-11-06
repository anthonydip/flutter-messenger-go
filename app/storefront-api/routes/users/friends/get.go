package friends

import (
	"encoding/json"
	"net/http"

	"github.com/anthonydip/flutter-messenger-go/app/storefront-api/utils"
	"github.com/anthonydip/flutter-messenger-go/app/storefront-api/webserver"
	"github.com/anthonydip/flutter-messenger-go/pkg/dtos"

	"github.com/rs/zerolog/log"
)

type FriendsResponse struct {
	Status        string        `json:"status"`
	StatusCode    int           `json:"statusCode"`
	StatusMessage string        `json:"statusMessage,omitempty"`
	Friends       []dtos.Friend `json:"friends"`
}

// Get friends list
func Get(srv webserver.Server) http.HandlerFunc {
	if srv == nil {
		log.Fatal().Msg("a nil dependency was passed to GET '/users/friends'")
	}

	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		log.Info().Msg("[GET /users/friends] Received a request")

		// Retrieve the access token
		authHeader := r.Header.Get("Authorization")
		token, err := utils.GetAuthorizationToken(authHeader)
		if err != nil {
			res := FriendsResponse{
				Status:     "UNAUTHORIZED",
				StatusCode: 401,
			}

			switch err.Error() {
			case "empty header":
				log.Error().Msgf("[GET /users/friends] Empty authorization header provided")
				res.StatusMessage = "Empty authorization header"
			case "invalid header":
				log.Error().Msgf("[GET /users/friends] Invalid authorization header provided")
				res.StatusMessage = "Invalid authorization header"
			default:
				res = FriendsResponse{
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

		user, err := srv.ValidateParseJWT(token)
		if err != nil {
			switch err.Error() {
			case "error reading pem":
				log.Error().Msgf("[GET /users/friends] Error reading PEM for token")
			case "error parsing pem":
				log.Error().Msgf("[GET /users/friends] Error parsing PEM for token")
			case "invalid token":
				res := FriendsResponse{
					Status:        "UNAUTHORIZED",
					StatusCode:    401,
					StatusMessage: "Invalid authorization token",
				}
				w.WriteHeader(http.StatusUnauthorized)
				json.NewEncoder(w).Encode(&res)
				return
			default:
				log.Error().Msgf("[GET /users/friends] Error occurred validating and parsing token")
			}

			res := FriendsResponse{
				Status:        "INTERNAL SERVER ERROR",
				StatusCode:    500,
				StatusMessage: "Error validating and parsing token",
			}
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(&res)
			return
		}

		sublogger := log.With().Any("user", user.Id).Logger()

		friendsList, err := srv.GetAllFriends(user.Id)
		if err != nil {
			sublogger.Info().Msgf("[GET /users/friends] Error getting friends from the database, %s", err.Error())

			res := FriendsResponse{
				Status:        "INTERNAL SERVER ERROR",
				StatusCode:    500,
				StatusMessage: "Error retrieving user friends list",
			}
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(&res)
			return
		}

		sublogger.Info().Msg("[GET /users/friends] Successfully retrieved user friends list")

		res := FriendsResponse{
			Status:        "SUCCESS",
			StatusCode:    200,
			StatusMessage: "Friends list retrieved",
			Friends:       friendsList,
		}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(&res)

	}
}
