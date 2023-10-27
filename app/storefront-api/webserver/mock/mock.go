package mock

import (
	"github.com/anthonydip/flutter-messenger-go/internal/storefront"
	"github.com/anthonydip/flutter-messenger-go/pkg/authentication"

	mockstore "github.com/anthonydip/flutter-messenger-go/internal/storefront/mock"
	// mockauth 	"github.com/anthonydip/flutter-messenger-go/pkg/authentication/mock"
)

type Result func(c *mockConfig)

type mockConfig struct {
}

// Mock the webserver
type Mock struct {
	authentication.Authentication
	storefront.Storefront

	cfg mockConfig
}

// Function to create a new Mock webserver
func New(opts ...Result) Mock {
	r := Mock{
		// Authentication: mockauth.New(),
		Storefront: mockstore.New(),
	}

	for _, o := range opts {
		if o != nil {
			o(&r.cfg)
		}
	}

	return r
}

// WithStorefront attaches a customized storefront mock
func (m Mock) WithStorefront(opts ...mockstore.Result) Mock {
	m.Storefront = mockstore.New(opts...)

	return m
}
