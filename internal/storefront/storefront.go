package storefront

import (
	"fmt"

	"github.com/anthonydip/flutter-messenger-go/pkg/dtos"
)

// Storefront exposes all functionalities of the Storefront service.
type Storefront interface {
	GetUser(string) (dtos.User, error)
}

// Broker manages the internal state of the Storefront service.
type Broker struct {
	cfg Config // the storefront's configuration
}

// New initializes a new Storefront service.
func New(cfg Config) (*Broker, error) {
	r := &Broker{}

	if err := validateConfig(cfg); err != nil {
		return nil, fmt.Errorf("invalid configuration: %w", err)
	}

	return r, nil
}
