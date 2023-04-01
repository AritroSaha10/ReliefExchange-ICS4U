package main

import (
	"context"
	"log"
	"net/http"

	"cloud.google.com/go/firestore"
	firebase "firebase.google.com/go"
	"github.com/gin-gonic/gin"
	"google.golang.org/api/iterator"
	"google.golang.org/api/option"
)

const SERVICE_ACCOUNT_FILENAME = "ics4u0-project-firebase-key.json"

var firebaseContext context.Context
var firebaseApp *firebase.App
var firestoreClient *firestore.Client

type Donation struct {
	ID          string   `json:"id"`
	Title       string   `json:"title"`
	Description string   `json:"description"`
	Location    string   `json:"location"`
	Images      []string `json:"imgs"`
}

func getDonationsListEndpoint(c *gin.Context) {
	donations, err := getAllDonations(firebaseContext, firestoreClient)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	} else {
		c.IndentedJSON(http.StatusOK, donations)
	}
}

func getDonationFromIDEndpoint(c *gin.Context) {
	id := c.Param("id")
	donation, err := getDonationByID(firebaseContext, firestoreClient, id)
	if err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"error": err.Error()})
	} else {
		c.IndentedJSON(http.StatusOK, donation)
	}
}

func postDonationEndpoint(c *gin.Context) {
	var donation Donation
	//set donation ownerId to current userId
	if err := c.ShouldBindJSON(&donation); err != nil { //transfers request body so that feilds match the donation struct
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	docID, err := addDonation(firebaseContext, firestoreClient, donation) //create new donation object from struct
	//add to the firestore databse
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	} else {
		c.IndentedJSON(http.StatusCreated, docID)
	}

}

func main() {
	firebaseContext = context.Background()
	firebaseCreds := option.WithCredentialsFile(SERVICE_ACCOUNT_FILENAME)

	app, err := firebase.NewApp(firebaseContext, nil, firebaseCreds)
	if err != nil {
		log.Fatalf("Error initializing Firebase app: %v\n", err)
	}
	firebaseApp = app

	firestoreClient, err = firebaseApp.Firestore(firebaseContext)
	if err != nil {
		log.Fatalf("Error initializing Firestore client: %v\n", err)
	}

	r := gin.Default()

	r.GET("/donations/donationList", getDonationsListEndpoint)
	r.GET("/donations/:id", getDonationFromIDEndpoint)
	r.POST("/donations/new", postDonationEndpoint)

	err = r.Run()
	if err != nil {
		return
	}
}

// Gets all donations available
func getAllDonations(ctx context.Context, client *firestore.Client) ([]Donation, error) {
	var donations []Donation
	iter := client.Collection("donations").Documents(ctx) //.Documents(ctx) returns a iterator
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, err //no data was retrieved-nil, but there was an error -err
		}
		var donation Donation
		err = doc.DataTo(&donation)
		if err != nil {
			return nil, err
		}
		donation.ID = doc.Ref.ID //sets donation struct id to the one in the firebase
		donations = append(donations, donation)
	}
	return donations, nil //nil-data was retrived without any errors
}

func getDonationByID(ctx context.Context, client *firestore.Client, id string) (Donation, error) {
	var donation Donation
	doc, err := client.Collection("donations").Doc(id).Get(ctx) //get a single donation from its id
	if err != nil {
		return donation, err //returns empty donation struct
	}
	err = doc.DataTo(&donation)
	if err != nil {
		return Donation{}, err
	}
	donation.ID = doc.Ref.ID
	return donation, nil
}

func addDonation(ctx context.Context, client *firestore.Client, donation Donation) (string, error) {
	docRef, _, err := client.Collection("donations").Add(ctx, map[string]interface{}{
		"title":       donation.Title,
		"description": donation.Description,
		"location":    donation.Location,
		"imgs":        donation.Images,
	})
	if err != nil {
		return "", err
	}
	return docRef.ID, nil
}
