package post

import (
	"net/http"
	"relief_exchange_backend/globals"
	"relief_exchange_backend/helpers"
	"relief_exchange_backend/types"

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
	//Bind the request body to the body struct, this stores the donation data and id token of the user to allow go to use.
	if err := c.ShouldBindJSON(&body); err != nil {
		log.Error(err.Error())
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	//verify the IdToken of the sender (user) with the server.
	token, err := globals.AuthClient.VerifyIDToken(globals.FirebaseContext, body.IDToken)
	if err != nil {
		log.Warn("Failed to verify ID token")
		c.IndentedJSON(http.StatusForbidden, gin.H{"error": "You are not authorized to make this donation."})
		return
	}

	userUID := token.UID
	//Use addDonation function to add the donation, passing in the donationData and the userid.
	docID, err := helpers.AddDonation(body.DonationData, userUID)
	//if there was an error adding the donation, send back the err message to the frontend, otherwise send back the docId for the frontend to use
	if err != nil {
		log.Error(err.Error())
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	} else {
		log.Info("Post donation successful.")
		c.IndentedJSON(http.StatusCreated, docID)
	}
}
