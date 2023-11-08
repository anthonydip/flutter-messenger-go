package webserver

import (
	"errors"
	"fmt"
	"net"
	"net/http"
	"os"
	"strconv"

	"github.com/anthonydip/flutter-messenger-go/app/storefront-api/ws"
	"github.com/anthonydip/flutter-messenger-go/internal/storefront"
	"github.com/anthonydip/flutter-messenger-go/pkg/authentication"

	"github.com/gorilla/mux"
	"github.com/rs/zerolog/log"
)

// Server exposes all functionalities of the Storefront API
type Server interface {
	authentication.Authentication
	storefront.Storefront
}

// Broker manages the internal state of the Storefront API
type Broker struct {
	authentication.Authentication
	storefront.Storefront

	cfg    Config      // the api service's configuration
	router *mux.Router // the api service's route collection
}

// Initialize a new Storefront API
func New(cfg Config) (*Broker, error) {
	r := &Broker{}

	err := validateConfig(cfg)
	if err != nil {
		return nil, fmt.Errorf("invalid configuration: %w", err)
	}
	r.cfg = cfg

	r.Authentication, err = authentication.New(cfg.Auth)
	if err != nil {
		return nil, fmt.Errorf("invalid auth configuration: %w", err)
	}

	r.Storefront, err = storefront.New(cfg.Storefront)
	if err != nil {
		return nil, fmt.Errorf("invalid storefront configuration: %w", err)
	}

	return r, nil
}

// Start the Storefront service
func (bkr *Broker) Start(binder func(s Server, h *ws.Hub, r *mux.Router)) {
	hub := ws.NewHub()
	go hub.Run()
	bkr.router = mux.NewRouter().StrictSlash(true)
	binder(bkr, hub, bkr.router)

	l, err := net.Listen("tcp", ":"+strconv.Itoa(bkr.cfg.Port))
	if err != nil {
		log.Fatal().Err(err).Msgf("Failed to bind to TCP port %d for listening.", bkr.cfg.Port)
		os.Exit(13)
	} else {
		log.Info().Msgf("Starting webserver on TCP port %04d", bkr.cfg.Port)
	}

	if err := http.Serve(l, bkr.router); errors.Is(err, http.ErrServerClosed) {
		log.Warn().Err(err).Msg("Web server has shut down")
	} else {
		log.Fatal().Err(err).Msg("Web server has shut down unexpectedly")
	}
}
