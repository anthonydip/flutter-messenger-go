package storefront

import (
	"context"
	"fmt"

	firebase "firebase.google.com/go"
)

func initializeFirebase(r *Broker) error {
	// Initialize firebase app
	conf := &firebase.Config{ProjectID: "***REMOVED***"}
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
