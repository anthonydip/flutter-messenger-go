package storefront

import (
	"context"

	"github.com/anthonydip/flutter-messenger-go/pkg/dtos"
)

type TokenInfo struct {
	TokenType string
	Email     string
	Provider  string
}

// Function to add the access token to the database
func (bkr Broker) AddAccessToken(token string, user dtos.User) error {
	info := TokenInfo{
		"user",
		user.Email,
		user.Provider,
	}

	_, err := bkr.Firestore.Collection("tokens").Doc(token).Set(context.Background(), info)
	if err != nil {
		return err
	}

	return nil
}
