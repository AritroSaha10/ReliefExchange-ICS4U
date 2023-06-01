package globals

import (
    "context"
    "os"

    "cloud.google.com/go/firestore"
    firebase "firebase.google.com/go"
    "firebase.google.com/go/auth"
    log "github.com/sirupsen/logrus"
    "google.golang.org/api/option"
)

// firebaseContext and firebaseApp are global variables for firebase context and application.
// firestoreClient and authClient are clients for firestore and authentication respectively.
var (
    FirebaseContext context.Context
    FirebaseApp     *firebase.App
    FirestoreClient *firestore.Client
    AuthClient      *auth.Client
)

func InitializeFirebaseGlobals() error {
    // Set up all Firebase connections
    // Set up context and import Firebase credentials
    FirebaseContext = context.Background()
    firebaseCreds := option.WithCredentialsJSON([]byte(os.Getenv("FIREBASE_CREDENTIALS_JSON")))

    // Set up Firebase
    FirebaseApp, err := firebase.NewApp(FirebaseContext, nil, firebaseCreds)
    if err != nil {
        log.Fatalf("Error initializing Firebase app: %v\n", err)
        return err
    }

    // Set up Firestore
    FirestoreClient, err = FirebaseApp.Firestore(FirebaseContext)
    if err != nil {
        log.Fatalf("Error initializing Firestore client: %v\n", err)
        return err
    }

    // Set up Firebase Auth
    AuthClient, err = FirebaseApp.Auth(FirebaseContext)
    if err != nil {
        log.Fatalf("Error initializing Firebase Auth client: %v\n", err)
        return err
    }

    return nil
}
