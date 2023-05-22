package globals

import (
	"context"

	"cloud.google.com/go/firestore"
	firebase "firebase.google.com/go"
	"firebase.google.com/go/auth"
)

// firebaseContext and firebaseApp are global variables for firebase context and application.
// firestoreClient and authClient are clients for firestore and authentication respectively.
var (
	FirebaseContext context.Context
	FirebaseApp     *firebase.App
	FirestoreClient *firestore.Client
	AuthClient      *auth.Client
)
