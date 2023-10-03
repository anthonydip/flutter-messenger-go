package authentication

import (
	"os"
	"time"

	"github.com/anthonydip/flutter-messenger-go/pkg/dtos"
	"github.com/golang-jwt/jwt"
)

type UserInfo struct {
	Email string
}

type JwtClaims struct {
	*jwt.StandardClaims
	TokenType string
	UserInfo
}

// Auth exposes all functionalities of the Auth agent
type Authentication interface {
	GenerateAccessToken(user dtos.User) (string, error)
	ValidateJWT(token string) bool
}

// Broker manages the internal state of the Auth agent.
type Broker struct{}

// New create a new authorization agent.
func New(cfg Config) (Authentication, error) {
	return &Broker{}, nil
}

func (bkr *Broker) GenerateAccessToken(user dtos.User) (string, error) {
	signBytes, err := os.ReadFile("../keys/rsa-access-key.private")
	if err != nil {
		return "", nil
	}

	signKey, err := jwt.ParseRSAPrivateKeyFromPEM(signBytes)
	if err != nil {
		return "", nil
	}

	// Create a signer for RSA 256
	t := jwt.New(jwt.GetSigningMethod("RS256"))

	// Set claims
	t.Claims = &JwtClaims{
		&jwt.StandardClaims{
			// Set the expire time
			// see http://tools.ietf.org/html/draft-ietf-oauth-json-web-token-20#section-4.1.4
			ExpiresAt: time.Now().Add(time.Minute * 1).Unix(),
		},
		"level1",
		UserInfo{
			Email: user.Email,
		},
	}

	// Create token string
	token, err := t.SignedString(signKey)
	if err != nil {
		return "", nil
	}

	return token, nil
}

// ValidateJWT validates the JWT token
func (bkr *Broker) ValidateJWT(token string) bool {
	// Do JWT validation stuff here
	return true
}
