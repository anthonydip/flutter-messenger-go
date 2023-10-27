package authentication

import (
	"fmt"
	"os"
	"time"

	"github.com/anthonydip/flutter-messenger-go/pkg/dtos"
	"github.com/golang-jwt/jwt/v5"
	"github.com/rs/zerolog/log"
)

const (
	privAccessKeyPath   = "../../keys/rsa-access-key.private"
	pubAccessKeyPath    = "../../keys/rsa-access-key.public"
	privInternalKeyPath = "../../keys/rsa-internal-key.private"
	pubInternalKeyPath  = "../../keys/rsa-internal-key.public"
)

type JwtClaims struct {
	jwt.RegisteredClaims
	TokenType string
	UserID    string
	Email     string
	Provider  string
}

// Auth exposes all functionalities of the Auth agent
type Authentication interface {
	GenerateAccessToken(user dtos.User) (string, error)
	ValidateJWT(token string) bool
	ValidateParseJWT(token string) (dtos.User, error)
	ValidateInternalJWT(token string) bool
}

// Broker manages the internal state of the Auth agent.
type Broker struct{}

// New create a new authorization agent.
func New(cfg Config) (Authentication, error) {
	return &Broker{}, nil
}

// Generate a new access token for a user
func (bkr *Broker) GenerateAccessToken(user dtos.User) (string, error) {
	// Read the private PEM key for the user access token
	signBytes, err := os.ReadFile(privAccessKeyPath)
	if err != nil {
		log.Fatal().Err(err).Str("function", "GenerateAccessToken").Msg("Error reading private PEM key")
		return "", fmt.Errorf("error reading file")
	}

	// Parse RSA from the private key
	signKey, err := jwt.ParseRSAPrivateKeyFromPEM(signBytes)
	if err != nil {
		log.Fatal().Err(err).Str("function", "GenerateAccessToken").Msg("Error parsing private PEM key")
		return "", fmt.Errorf("error parsing pem")
	}

	// Create claims with a 1 minute expire time
	claims := JwtClaims{
		jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
		},
		"user",
		user.Id,
		user.Email,
		user.Provider,
	}

	// Create new token with claims and sign
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	ss, err := token.SignedString(signKey)
	if err != nil {
		log.Fatal().Err(err).Str("function", "GenerateAccessToken").Msg("Error signing token")
		return "", fmt.Errorf("error signing token")
	}

	return ss, nil
}

// ValidateJWT validates the user JWT token
func (bkr *Broker) ValidateJWT(tokenString string) bool {
	// Read the public PEM key for the user access token
	verifyBytes, err := os.ReadFile(pubAccessKeyPath)
	if err != nil {
		log.Fatal().Err(err).Str("function", "ValidateJWT").Msg("Error reading public PEM key")
		return false
	}

	// Parse RSA from the public key
	verifyKey, err := jwt.ParseRSAPublicKeyFromPEM(verifyBytes)
	if err != nil {
		log.Fatal().Err(err).Str("function", "ValidateJWT").Msg("Error parsing public PEM key")
		return false
	}

	// Verify the provided token string
	_, err = jwt.ParseWithClaims(tokenString, &JwtClaims{}, func(token *jwt.Token) (interface{}, error) {
		return verifyKey, nil
	})

	// If the token is missing or invalid, return false
	return err == nil
}

// Function to validate user JWT token and return the associated user information
func (bkr *Broker) ValidateParseJWT(tokenString string) (dtos.User, error) {
	// Read the public PEM key for the user access token
	verifyBytes, err := os.ReadFile(pubAccessKeyPath)
	if err != nil {
		return dtos.User{}, fmt.Errorf("error reading pem")
	}

	// Parse RSA from the public key
	verifyKey, err := jwt.ParseRSAPublicKeyFromPEM(verifyBytes)
	if err != nil {
		return dtos.User{}, fmt.Errorf("error parsing pem")
	}

	// Verify the provided token string
	token, err := jwt.ParseWithClaims(tokenString, &JwtClaims{}, func(token *jwt.Token) (interface{}, error) {
		return verifyKey, nil
	})

	if claims, ok := token.Claims.(*JwtClaims); ok && token.Valid {
		return dtos.User{
			Id:       claims.UserID,
			Email:    claims.Email,
			Provider: claims.Provider,
		}, nil
	} else {
		return dtos.User{}, fmt.Errorf("invalid token")
	}
}

// ValidateInternalJWT validates the internal JWT token
func (bkr *Broker) ValidateInternalJWT(tokenString string) bool {
	// Read the public PEM key for the internal access token
	verifyBytes, err := os.ReadFile(pubInternalKeyPath)
	if err != nil {
		log.Fatal().Err(err).Str("function", "ValidateInternalJWT").Msg("Error reading public PEM key")
		return false
	}

	// Parse RSA from the public key
	verifyKey, err := jwt.ParseRSAPublicKeyFromPEM(verifyBytes)
	if err != nil {
		log.Fatal().Err(err).Str("function", "ValidateInternalJWT").Msg("Error parsing public PEM key")
		return false
	}

	// Verify the provided token string
	_, err = jwt.ParseWithClaims(tokenString, &JwtClaims{}, func(token *jwt.Token) (interface{}, error) {
		return verifyKey, nil
	})

	// If the token is missing or invalid, return false
	return err == nil
}
