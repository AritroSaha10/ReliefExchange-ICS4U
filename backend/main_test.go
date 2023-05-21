package main

import (
	"context"
	"log"
	"math/rand"
	"os"
	"time"

	"testing"

	"github.com/stretchr/testify/assert"

	"cloud.google.com/go/firestore"
	firebase "firebase.google.com/go"
	"github.com/getsentry/sentry-go"
	_ "github.com/joho/godotenv/autoload"

	"google.golang.org/api/option"
)

func TestMain(m *testing.M) {
	firebaseContext = context.Background()
	firebaseCreds := option.WithCredentialsJSON([]byte(os.Getenv("FIREBASE_CREDENTIALS_JSON")))

	err := sentry.Init(sentry.ClientOptions{
		Dsn:              os.Getenv("https://4044f25736934d42862ea077a1283931@o924596.ingest.sentry.io/4505213654073344"),
		TracesSampleRate: 1.0,
	})
	if err != nil {
		log.Fatalf("Error initializing Sentry: %s", err)
	}

	app, err := firebase.NewApp(firebaseContext, nil, firebaseCreds)
	if err != nil {
		log.Fatalf("Error initializing Firebase app: %v\n", err)
	}
	firebaseApp = app

	firestoreClient, err = firebaseApp.Firestore(firebaseContext)
	if err != nil {
		log.Fatalf("Error initializing Firestore client: %v\n", err)
	}

	authClient, err = firebaseApp.Auth(firebaseContext)
	if err != nil {
		log.Fatalf("Error initializing Firebase Auth client: %v\n", err)
	}

	if err != nil {
		return
	}
	os.Exit(m.Run())
}

func TestGetAllDonations(t *testing.T) {
	donations, err := getAllDonations(firebaseContext, firestoreClient)
	assert.NoError(t, err, "getAllDonations function should return without error")
	assert.NotEmpty(t, donations, "getAllDonations should return at least one donation")

}
func TestAddDonation(t *testing.T) {
	donation := Donation{
		ID:                "testID",
		Title:             "testTitle",
		Description:       "testDescription",
		Location:          "testLocation",
		Image:             "testImage",
		CreationTimestamp: time.Now().UTC(),
		OwnerId:           "testOwnerId",
		Tags:              []string{"tag1", "tag2"},
		Reports:           []string{"report1", "report2"},
	}
	test_user_id := "4P9lIlcIYNeeCZja6Wc3coemX1A3" //Joshua.C
	donationId, err := addDonation(firebaseContext, firestoreClient, donation, test_user_id)
	assert.NoError(t, err, "addDonation function should return without error")
	assert.NotEmpty(t, donationId, "addDonation should return a donation id ")
	owner, err := firestoreClient.Doc("users/" + test_user_id).Get(firebaseContext)

	rawPosts := owner.Data()["posts"].([]interface{})
	var posts []string
	for _, value := range rawPosts {
		docRef, ok := value.(*firestore.DocumentRef)
		if !ok {
			t.Errorf("Error: value is not a *firestore.DocumentRef")
			continue
		}
		posts = append(posts, docRef.ID)
	}
	assert.Contains(t, posts, donationId, "add Donation should add the donation to the user posts feild")

}
func TestGetDonationById(t *testing.T) {
	test_user_id := "4P9lIlcIYNeeCZja6Wc3coemX1A3" //Joshua.C
	owner, err := firestoreClient.Doc("users/" + test_user_id).Get(firebaseContext)
	rawPosts := owner.Data()["posts"].([]interface{})
	var posts []string
	for _, value := range rawPosts {
		docRef, ok := value.(*firestore.DocumentRef)
		if !ok {
			t.Errorf("Error: value is not a *firestore.DocumentRef")
			continue
		}
		posts = append(posts, docRef.ID)
	}
	donation, err := getDonationByID(firebaseContext, firestoreClient, posts[rand.Intn(len(posts)-1)])
	assert.NoError(t, err, "GetDonationById function should return without error")
	assert.NotEmpty(t, donation, "GetDonationsById should return at least one donation")
	assert.False(t, donation.CreationTimestamp.IsZero(), "CreationTimestamp should be set")

}
func TestCheckIfAdmin(t *testing.T) {
	test_user_id := "4P9lIlcIYNeeCZja6Wc3coemX1A3" //Joshua.C
	isAdmin, err := checkIfAdmin(firebaseContext, firestoreClient, test_user_id)
	assert.NoError(t, err, "GetDonationById function should return without error")
	assert.False(t, isAdmin, "Joshua.C is not an admin")
}
