package accessToken

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/anthonydip/flutter-messenger-go/app/storefront-api/webserver"
	"github.com/gorilla/mux"
	"github.com/rs/zerolog/log"
)

type DeleteResponse struct {
	Status        string `json:"status"`
	StatusCode    int    `json:"statusCode"`
	StatusMessage string `json:"statusMessage,omitempty"`
}

// Delete a user access token
func Delete(srv webserver.Server) http.HandlerFunc {
	if srv == nil {
		log.Fatal().Msg("a nil dependency was passed to DELETE '/tokens/access/{token}'")
	}

	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		// Get the token
		params := mux.Vars(r)
		token := strings.TrimSpace(params["token"])

		sublogger := log.With().Any("token", token).Logger()
		sublogger.Info().Msg("[DELETE /auth/tokens/{token}] Received a request")

		// Check if the token exists
		err := srv.AccessTokenExists(token)
		if err != nil {
			if err.Error() == "token not found" {
				sublogger.Error().Msg("[DELETE /auth/tokens/{token}] Token not found")

				res := DeleteResponse{
					Status:        "NOT FOUND",
					StatusCode:    404,
					StatusMessage: "Token does not exist",
				}
				w.WriteHeader(http.StatusNotFound)
				json.NewEncoder(w).Encode(&res)
				return
			} else {
				sublogger.Error().Msgf("[DELETE /auth/tokens/{token}] Error checking for token, %s", err.Error())

				res := DeleteResponse{
					Status:        "INTERNAL SERVER ERROR",
					StatusCode:    500,
					StatusMessage: "Error checking for token",
				}
				w.WriteHeader(http.StatusInternalServerError)
				json.NewEncoder(w).Encode(&res)
				return
			}
		}

		// Delete the token from the database
		err = srv.DeleteAccessToken(token)
		if err != nil {
			sublogger.Error().Msgf("[DELETE /auth/tokens/{token}] Error deleting token from database, %s", err.Error())

			res := DeleteResponse{
				Status:        "INTERNAL SERVER ERROR",
				StatusCode:    500,
				StatusMessage: "Error deleting token",
			}
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(&res)
			return
		}

		sublogger.Info().Msg("[DELETE /auth/tokens/{token}] Token successfully deleted")

		res := DeleteResponse{
			Status:        "SUCCESS",
			StatusCode:    200,
			StatusMessage: "Token successfully deleted",
		}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(&res)
	}
}
