/*
Package post contains endpoints for handling actions related to making donations.
This includes adding new donations to the database, and verifying the authenticity
of the donation user through token verification.

This package includes the following import dependencies:

    "net/http" : Provides HTTP client and server implementations
    "relief_exchange_backend/globals" : Contains global variables or objects
    "relief_exchange_backend/helpers" : Contains helper functions
    "relief_exchange_backend/types" : Contains types that are used in the backend
    "github.com/gin-gonic/gin" : Gin is a HTTP web framework written in Go
    "github.com/sirupsen/logrus" : Logrus is a structured logger for Go

The AddDonation function handles the endpoint to post a new donation. The function
verifies the token of the user making the donation and adds the donation to the database.
If the donation is added successfully, the function returns the document ID of the donation,
otherwise it returns an error.
// @authors Joshua Chou,Aritro Saha
@cite "Package iter." Pkg.go.dev, 2023. [Online].
Available: https://pkg.go.dev/github.com/reiver/go-iter. [Accessed: 30- May- 2023].

*/

package post

import (
	"net/http"
	"relief_exchange_backend/globals"
	"relief_exchange_backend/helpers"
	"relief_exchange_backend/types"

	"cloud.google.com/go/firestore"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

// postDonationEndpoint handles the endpoint to post a new donation.
// Parameters:
//   - c: the gin context, the request and response http.
//
// It accepts a donation and a user's id token, verifies the token,
// and then uses the addDonation function to add the donation to the database.
func AddDonation(c *gin.Context) {
	var body struct {
		DonationData types.Donation `json:"data"`
		IDToken      string         `json:"token"`
	}
	// Bind the request body to the body struct, this stores the donation data and id token of the user to allow go to use.
	if err := c.ShouldBindJSON(&body); err != nil {
		log.Error(err.Error())
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	// Verify the IdToken of the sender (user) with the server
	token, err := globals.AuthClient.VerifyIDToken(globals.FirebaseContext, body.IDToken)
	if err != nil {
		log.Warn("Failed to verify ID token")
		c.IndentedJSON(http.StatusForbidden, gin.H{"error": "You are not authorized to make this donation."})
		return
	}

	userUID := token.UID

	// Use addDonation function to add the donation, passing in the donationData and the uid
	docID, err := helpers.AddDonation(body.DonationData, userUID)

	// If there's an error adding the donation, send back err msg to frontend,
	// otherwise send back docId for the frontend to use
	if err != nil {
		log.Error(err.Error())
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	} else {
		log.Info("Post donation successful.")
		c.IndentedJSON(http.StatusCreated, docID)
	}

	// Increment the user's donation counter.
	// Run this after setting up the response since this isn't a priority
	// and won't affect it.
	userData, err := helpers.GetUserDataByID(token.UID)
	// Only do it if there was no error
	// Don't bother returning an error to endpoint since updating the counter
	// isn't a proper failure
	if err == nil {
		// Update user data with new donations count
		_, err = globals.FirestoreClient.Collection("users").Doc(token.UID).Update(globals.FirebaseContext, []firestore.Update{
			{
				Path:  "donations_made",
				Value: userData.DonationsMade + 1,
			},
		})
	}

	// Log any errors that occured
	if err != nil {
		log.Error(err.Error())
	}
}
