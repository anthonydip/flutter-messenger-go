package webserver

import (
	"errors"
	"net"
	"net/http"
	"os"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/rs/zerolog/log"
)

// Server exposes all functionalitys of the Messenger API
type Server interface {
	// TO-DO
}

// Broker manages the internal state of the Messenger API
type Broker struct {
	port   int
	router *mux.Router
}

// Initialize broker + validate configurations and run pre-flight checks
func New(port int) (*Broker, error) {
	r := &Broker{
		port: port,
	}

	// Perform any checks if needed

	return r, nil
}

// Start the Messenger service
func (bkr *Broker) Start(binder func(s Server, r *mux.Router)) {
	bkr.router = mux.NewRouter().StrictSlash(true)
	binder(bkr, bkr.router)

	l, err := net.Listen("tcp", ":"+strconv.Itoa(bkr.port))
	if err != nil {
		log.Fatal().Err(err).Msgf("Failed to bind to TCP port %d for listening.", bkr.port)
		os.Exit(13)
	} else {
		log.Info().Msgf("Starting webserver on TCP port %04d", bkr.port)
	}

	if err := http.Serve(l, bkr.router); errors.Is(err, http.ErrServerClosed) {
		log.Warn().Err(err).Msg("Web server has shut down")
	} else {
		log.Fatal().Err(err).Msg("Web server has shut down unexpectedly")
	}
}
