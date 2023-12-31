package mock

import (
	"github.com/anthonydip/flutter-messenger-go/pkg/dtos"
)

type Result func(c *mockConfig)

type mockConfig struct {
	validateJWTShouldFail bool
}

// Mock the Authorization agent
type Mock struct {
	cfg mockConfig
}

// Function to create a new Mock Authorization agent
func New(opts ...Result) *Mock {
	r := &Mock{}

	for _, o := range opts {
		if o != nil {
			o(&r.cfg)
		}
	}

	return r
}

func (m Mock) GenerateAccessToken(dtos.User) (string, error) {
	return "some-access-token", nil
}

func (m Mock) ValidateJWT(string) bool {
	return !m.cfg.validateJWTShouldFail
}

func ValidateJWTFail() Result {
	return func(c *mockConfig) {
		c.validateJWTShouldFail = true
	}
}

func (m Mock) ValidateParseJWT(string) (dtos.User, error) {
	return dtos.User{}, nil
}

func (m Mock) ValidateInternalJWT(string) bool {
	return true
}
