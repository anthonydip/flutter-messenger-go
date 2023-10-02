package storefront

import (
	"context"
	"fmt"
	"strconv"

	"github.com/anthonydip/flutter-messenger-go/pkg/dtos"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/api/iterator"
)

func (bkr Broker) GetUser(id string) (dtos.User, error) {
	dsnap, err := bkr.Firestore.Collection("users").Doc(id).Get(context.Background())
	if err != nil {
		fmt.Printf("An error has occurred: %s\n", err)
	}

	m := dsnap.Data()
	fmt.Printf("document data: %#v\n", m)

	return dtos.User{
		Email: "test@gmail.com",
	}, nil
}

func (bkr Broker) PostUser(userInfo dtos.User) (dtos.User, error) {
	// Check if the email already exists in firestore
	iter := bkr.Firestore.Collection("users").Where("email", "==", userInfo.Email).Limit(1).Documents(context.Background())
	for {
		_, err := iter.Next()
		if err == iterator.Done {
			break
		}
		// Error querying database
		if err != nil {
			return dtos.User{}, err
		}

		// Email already exists in the database
		return dtos.User{}, fmt.Errorf("409 Conflict")
	}

	// Get the salt rounds
	env, err := getEnv("SALT_ROUNDS")
	if err != nil {
		return dtos.User{}, err
	}

	// Convert salt rounds to a number
	saltRounds, err := strconv.Atoi(env)
	if err != nil {
		return dtos.User{}, err
	}

	// Hash the password
	hash, err := bcrypt.GenerateFromPassword([]byte(userInfo.Password), saltRounds)
	if err != nil {
		return dtos.User{}, err
	}

	// Generate a UUID for the new user
	id := uuid.New().String()

	user := dtos.User{
		Id:       id,
		Email:    userInfo.Email,
		Provider: userInfo.Provider,
		Password: string(hash),
	}

	// Add the new user to the database
	_, err = bkr.Firestore.Collection("users").Doc(id).Set(context.Background(), user)
	// Error creating user in the database
	if err != nil {
		return dtos.User{}, err
	}

	return user, nil
}
