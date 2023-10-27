package storefront

import (
	"context"
	"fmt"

	"github.com/anthonydip/flutter-messenger-go/pkg/dtos"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type TokenInfo struct {
	TokenType string
	Email     string
	Provider  string
}

// Function to check if an access token exists in the database
func (bkr Broker) AccessTokenExists(token string) error {
	_, err := bkr.Firestore.Collection("tokens").Doc(token).Get(context.Background())
	if err != nil {
		if status.Code(err) == codes.NotFound {
			return fmt.Errorf("token not found")
		}
		return err
	}

	return nil
}

// Function to delete the access token from the database
func (bkr Broker) DeleteAccessToken(token string) error {
	_, err := bkr.Firestore.Collection("tokens").Doc(token).Delete(context.Background())
	if err != nil {
		return err
	}

	return nil
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
