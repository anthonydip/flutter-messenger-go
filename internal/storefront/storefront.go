package storefront

import (
	"fmt"

	"github.com/anthonydip/flutter-messenger-go/pkg/dtos"

	"cloud.google.com/go/firestore"
)

// Storefront exposes all functionalities of the Storefront service.
type Storefront interface {
	GetUser(string) (dtos.User, error)
	GetUserByEmail(string) (dtos.User, error)
	PostUser(dtos.User) (dtos.User, error)
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

	err := initializeFirebase(r)
	if err != nil {
		return nil, err
	}

	return r, nil
}
