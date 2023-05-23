package main

import (
	"log"
	"math/rand"
	"os"
	"relief_exchange_backend/globals"
	"relief_exchange_backend/helpers"
	"relief_exchange_backend/types"
	"time"

	"testing"

	"github.com/stretchr/testify/assert"

	"cloud.google.com/go/firestore"
	"github.com/getsentry/sentry-go"
	_ "github.com/joho/godotenv/autoload"
)

func TestMain(m *testing.M) {
	err := sentry.Init(sentry.ClientOptions{
		Dsn:              "https://4044f25736934d42862ea077a1283931@o924596.ingest.sentry.io/4505213654073344",
		TracesSampleRate: 1.0,
	})
	if err != nil {
		log.Fatalf("Error initializing Sentry: %s", err)
	}

	// Initialize Firebase globals
	err = globals.InitializeFirebaseGlobals()
	if err != nil {
		return
	}
	os.Exit(m.Run())
}

func TestGetAllDonations(t *testing.T) {
	donations, err := helpers.GetAllDonations()
	assert.NoError(t, err, "getAllDonations function should return without error")
	assert.NotEmpty(t, donations, "getAllDonations should return at least one donation")
}

func TestAddDonation(t *testing.T) {
	donation := types.Donation{
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
	donationId, err := helpers.AddDonation(donation, test_user_id)
	assert.NoError(t, err, "addDonation function should return without error")
	assert.NotEmpty(t, donationId, "addDonation should return a donation id ")

	owner, err := globals.FirestoreClient.Doc("users/" + test_user_id).Get(globals.FirebaseContext)
	assert.NoError(t, err, "Owner should have been retrieved properly")

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
	owner, err := globals.FirestoreClient.Doc("users/" + test_user_id).Get(globals.FirebaseContext)
	assert.NoError(t, err, "Owner should have been retrieved properly")

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

	donation, err := helpers.GetDonationByID(posts[rand.Intn(len(posts)-1)])
	assert.NoError(t, err, "GetDonationById function should return without error")
	assert.NotEmpty(t, donation, "GetDonationsById should return at least one donation")
	assert.False(t, donation.CreationTimestamp.IsZero(), "CreationTimestamp should be set")
}

func TestCheckIfAdmin(t *testing.T) {
	test_user_id := "4P9lIlcIYNeeCZja6Wc3coemX1A3" //Joshua.C
	isAdmin, err := helpers.CheckIfAdmin(test_user_id)
	assert.NoError(t, err, "GetDonationById function should return without error")
	assert.False(t, isAdmin, "Joshua.C is not an admin")
}
