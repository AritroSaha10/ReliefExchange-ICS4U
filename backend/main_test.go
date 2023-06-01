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

	mockfs "github.com/weathersource/go-mockfs"

	pb "google.golang.org/genproto/googleapis/firestore/v1"
	"google.golang.org/protobuf/types/known/timestamppb"
)

var setup bool

func TestMain(m *testing.M) {
	if !setup {
		err := sentry.Init(sentry.ClientOptions{
			Dsn:              "https://4044f25736934d42862ea077a1283931@o924596.ingest.sentry.io/4505213654073344",
			TracesSampleRate: 1.0,
		})
		if err != nil {
			log.Fatalf("Error initializing Sentry: %s", err)
		}
		var t *testing.T
		client, server, err := mockfs.New()
		assert.NotNil(t, client)
		assert.NotNil(t, server)
		assert.Nil(t, err)

		globals.FirestoreClient = client

		// Initialize Firebase globals
		// Populate a mock document "b" in collection "C"
		var (
			aTime       = time.Date(2017, 1, 26, 0, 0, 0, 0, time.UTC)
			aTime2      = time.Date(2017, 2, 5, 0, 0, 0, 0, time.UTC)
			aTimestamp  = timestamppb.New(aTime)
			aTimestamp2 = timestamppb.New(aTime2)
			dbPath      = "projects/projectID/databases/(default)"
			// collectionPath        = "projects/projectID/databases/(default)/documents/donations"
			docPath = "projects/projectID/databases/(default)/documents/donations/interesting"
			pdoc    = &pb.Document{
				Name:       docPath,
				CreateTime: aTimestamp,
				UpdateTime: aTimestamp,
				Fields: map[string]*pb.Value{
					"id":                 {ValueType: &pb.Value_StringValue{"interesting"}},
					"title":              {ValueType: &pb.Value_StringValue{"interesting"}},
					"description":        {ValueType: &pb.Value_StringValue{"interesting"}},
					"location":           {ValueType: &pb.Value_StringValue{"interesting"}},
					"img":                {ValueType: &pb.Value_StringValue{"interesting"}},
					"creation_timestamp": {ValueType: &pb.Value_TimestampValue{aTimestamp}},
					"owner_id":           {ValueType: &pb.Value_StringValue{"interesting"}},
					"tags":               {ValueType: &pb.Value_ArrayValue{&pb.ArrayValue{}}},
					"reports":            {ValueType: &pb.Value_ArrayValue{&pb.ArrayValue{}}},
				},
			}
		)
		server.AddRPC(
			&pb.BatchGetDocumentsRequest{
				Database:  dbPath,
				Documents: []string{docPath},
			},
			[]interface{}{
				&pb.BatchGetDocumentsResponse{
					Result:   &pb.BatchGetDocumentsResponse_Found{pdoc},
					ReadTime: aTimestamp2,
				},
			},
		)

		setup = true
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
		Image:             "",
		CreationTimestamp: time.Now().UTC(),
		OwnerId:           "testOwnerId",
		Tags:              []string{"tag1", "tag2"},
		Reports:           []string{"report1", "report2"},
	}
	test_user_id := "p48oQ0SAYPeqculMRp2UBNJl03d2" //Joshua.C
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
	test_user_id := "p48oQ0SAYPeqculMRp2UBNJl03d2" //Joshua.C
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
	test_user_id := "p48oQ0SAYPeqculMRp2UBNJl03d2" //Joshua.C
	isAdmin, err := helpers.CheckIfAdmin(test_user_id)
	assert.NoError(t, err, "GetDonationById function should return without error")
	assert.True(t, isAdmin, "Joshua.C is an admin")
}
