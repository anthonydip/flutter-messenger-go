package utils

import (
	"fmt"
	"net/mail"

	"github.com/anthonydip/flutter-messenger-go/pkg/dtos"
	"github.com/caitlin615/nist-password-validator/password"
)

// Function to validate the fields from a create user request
func ValidatePostUser(user dtos.User) error {
	// Validate the email address
	_, err := mail.ParseAddress(user.Email)
	if err != nil {
		return fmt.Errorf("invalid email")
	}

	// Validate the provider
	if user.Provider != "Google" && user.Provider != "Flutter" {
		return fmt.Errorf("invalid provider")
	}

	// Validate the password
	validator := password.NewValidator(true, 8, 64)
	err = validator.ValidatePassword(user.Password)
	if err != nil {
		return fmt.Errorf("invalid password")
	}

	return nil
}
