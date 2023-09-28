package users

import (
	"fmt"
	"net/http"

	"github.com/anthonydip/flutter-messenger-go/app/messenger-api/webserver"

	"github.com/rs/zerolog/log"
)

func Get(srv webserver.Server) http.HandlerFunc {
	if srv == nil {
		log.Fatal().Msg("a nil dependency was passed to the '/users/{userID}'")
	}

	return func(w http.ResponseWriter, r *http.Request) {

		w.WriteHeader(http.StatusOK)
		w.Write([]byte(fmt.Sprintf("got user: %s", "user 1")))
	}
}
