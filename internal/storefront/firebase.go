package storefront

import (
	"context"
	"fmt"

	firebase "firebase.google.com/go"
)

func initializeFirebase(r *Broker) error {
	// Get Firebase Project ID from environment variable
	projectID, err := getEnv("FIREBASE_PROJECT_ID")
	if err != nil {
		return fmt.Errorf("error retrieving project id: %w", err)
	}

	// Initialize firebase app
	conf := &firebase.Config{ProjectID: projectID}
	app, err := firebase.NewApp(context.Background(), conf)
	if err != nil {
		return fmt.Errorf("error initializing firebase app: %w", err)
	}

	// Initialize firestore
	r.Firestore, err = app.Firestore(context.Background())
	if err != nil {
		return fmt.Errorf("error initializing firestore: %w", err)
	}

	return nil
}
