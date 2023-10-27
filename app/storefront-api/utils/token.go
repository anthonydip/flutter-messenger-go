package utils

import (
	"fmt"
	"strings"
)

// Function to extract the token from the authorization header
func GetAuthorizationToken(authHeader string) (string, error) {
	// Check for empty authorization header
	if authHeader == "" {
		return "", fmt.Errorf("empty header")
	}

	parts := strings.SplitN(authHeader, " ", 2)
	if len(parts) != 2 || parts[0] != "Bearer" {
		return "", fmt.Errorf("invalid header")
	}

	token := parts[1]

	return token, nil
}
