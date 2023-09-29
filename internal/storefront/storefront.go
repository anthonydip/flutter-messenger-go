package storefront

import (
	"context"
	"fmt"

	"github.com/anthonydip/flutter-messenger-go/pkg/dtos"

	"cloud.google.com/go/firestore"
	firebase "firebase.google.com/go"
)

// Storefront exposes all functionalities of the Storefront service.
type Storefront interface {
	GetUser(string) (dtos.User, error)
}

// Broker manages the internal state of the Storefront service.
type Broker struct {
	Firestore *firestore.Client

	cfg Config // the storefront's configuration
}

// New initializes a new Storefront service.
func New(cfg Config) (*Broker, error) {
	r := &Broker{}

	if err := validateConfig(cfg); err != nil {
		return nil, fmt.Errorf("invalid configuration: %w", err)
	}

	// Initialize firebase app
	conf := &firebase.Config{ProjectID: "***REMOVED***"}
	app, err := firebase.NewApp(context.Background(), conf)
	if err != nil {
		return nil, fmt.Errorf("error initializing firebase app: %w", err)
	}

	// Initialize firestore
	r.Firestore, err = app.Firestore(context.Background())
	if err != nil {
		return nil, fmt.Errorf("error initializing firestore: %w", err)
	}

	return r, nil
}
