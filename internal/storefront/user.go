package storefront

import (
	"context"
	"fmt"
	"strconv"

	"github.com/anthonydip/flutter-messenger-go/pkg/dtos"

	"github.com/google/uuid"
	"github.com/mitchellh/mapstructure"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/api/iterator"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (bkr Broker) GetUser(id string) (dtos.User, error) {
	user := dtos.User{}

	dsnap, err := bkr.Firestore.Collection("users").Doc(id).Get(context.Background())
	if err != nil {
		if status.Code(err) == codes.NotFound {
			return dtos.User{}, fmt.Errorf("user not found")
		}
		return dtos.User{}, err
	}

	mapstructure.Decode(dsnap.Data(), &user)

	user.Password = ""

	return user, nil
}

func (bkr Broker) GetUserByEmail(email string) (dtos.User, error) {
	user := dtos.User{}

	iter := bkr.Firestore.Collection("users").Where("email", "==", email).Limit(1).Documents(context.Background())
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}

		if err != nil {
			return dtos.User{}, err
		}

		// User exists in the database
		if doc.Data() != nil {
			mapstructure.Decode(doc.Data(), &user)
		}
	}

	// User does not exist
	if (dtos.User{}) == user {
		return dtos.User{}, fmt.Errorf("user does not exist")
	}

	user.Password = ""

	return user, nil
}

func (bkr Broker) GetAllFriends(id string) ([]dtos.Friend, error) {
	friends := make([]dtos.Friend, 0)

	iter := bkr.Firestore.Collection("users").Doc(id).Collection("friends").Documents(context.Background())
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}

		if err != nil {
			return make([]dtos.Friend, 0), err
		}

		// User exists in the database
		if doc.Data() != nil {
			friend := dtos.Friend{}

			mapstructure.Decode(doc.Data(), &friend)
			friends = append(friends, friend)
		}
	}

	return friends, nil
}

// Sign a user in, checking password
func (bkr Broker) SignIn(userInfo dtos.User) error {
	user := dtos.User{}

	// Search for user
	iter := bkr.Firestore.Collection("users").Where("email", "==", userInfo.Email).Limit(1).Documents(context.Background())
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}

		if err != nil {
			return err
		}

		// User exists in the database
		if doc.Data() != nil {
			mapstructure.Decode(doc.Data(), &user)
		}
	}

	// User does not exist
	if (dtos.User{}) == user {
		return fmt.Errorf("user does not exist")
	}

	// Validate the request password is the same as the hash password
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(userInfo.Password))
	if err != nil {
		return fmt.Errorf("invalid password")
	}

	return nil
}

func (bkr Broker) PostUser(userInfo dtos.User) (dtos.User, error) {
	user := dtos.User{}

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

	if userInfo.Provider == "Flutter" {
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

		user = dtos.User{
			Id:       id,
			Email:    userInfo.Email,
			Provider: userInfo.Provider,
			Password: string(hash),
		}
	} else {
		// Generate a UUID for the new user
		id := uuid.New().String()

		user = dtos.User{
			Id:       id,
			Email:    userInfo.Email,
			Provider: userInfo.Provider,
			Password: "",
		}
	}

	// Add the new user to the database
	_, err := bkr.Firestore.Collection("users").Doc(user.Id).Set(context.Background(), user)
	// Error creating user in the database
	if err != nil {
		return dtos.User{}, err
	}

	return user, nil
}

func (bkr Broker) PostFriend(userID string, friend dtos.User) error {
	// Check if the friend is already added
	_, err := bkr.Firestore.Collection("users").Doc(userID).Collection("friends").Doc(friend.Id).Get(context.Background())
	// Error is expected in this case
	if err != nil {
		// If friend is not added, add them
		if status.Code(err) == codes.NotFound {
			_, err = bkr.Firestore.Collection("users").Doc(userID).Collection("friends").Doc(friend.Id).Set(context.Background(), friend)
			if err != nil {
				return err
			}

			// Friend successfully added
			return nil
		} else {
			return err
		}
	}

	return fmt.Errorf("friend already added")
}
