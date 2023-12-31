package utils

import (
	"fmt"
	"net/mail"

	"github.com/anthonydip/flutter-messenger-go/pkg/dtos"
	"github.com/caitlin615/nist-password-validator/password"
)

// Function to validate an email
func ValidateEmail(email string) error {
	_, err := mail.ParseAddress(email)
	if err != nil {
		return fmt.Errorf("invalid email")
	}

	return nil
}

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

	// Validate the password if it is not through Google
	if user.Provider == "Flutter" {
		validator := password.NewValidator(true, 8, 64)
		err = validator.ValidatePassword(user.Password)
		if err != nil {
			return fmt.Errorf("invalid password")
		}
	}

	return nil
}
