package users

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/anthonydip/flutter-messenger-go/app/storefront-api/webserver"

	"github.com/gorilla/mux"
	"github.com/rs/zerolog/log"
)

func Get(srv webserver.Server) http.HandlerFunc {
	if srv == nil {
		log.Fatal().Msg("a nil dependency was passed to GET '/users/{userID}'")
	}

	return func(w http.ResponseWriter, r *http.Request) {
		params := mux.Vars(r)
		userID := strings.TrimSpace(params["userID"])

		user, err := srv.GetUser(userID)
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte("404 Not Found"))
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write([]byte(fmt.Sprintf("got user: %s", user.Name)))
	}
}
