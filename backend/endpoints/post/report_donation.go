/*
 * File: report_donation.go
 * -------------
 * This module handles the report donation endpoint in the server.
 * It takes a gin context as a parameter, binds the request body to a struct,
 * extracts the token and donation id from it, and verifies the token.
 * If the token is valid, it calls the ReportDonation helper function to report the donation with the given donation id.
 // @author Aritro Saha
*/

package post

import (
	"net/http"
	"relief_exchange_backend/globals"
	"relief_exchange_backend/helpers"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

// ReportDonation handles the endpoint to report a donation.
// Parameters:
//   - c: the gin context, the request and response http.
//
// It accepts a user's id token and a donation id, verifies the token,
// then adds the user's report to the donation.
func ReportDonation(c *gin.Context) {
	// Define body to store request information
	var body struct {
		DonationID string `json:"donation_id"`
		IDToken    string `json:"token"`
	}
	// Attempt to bind the request to the body, so golang can use the donation_id and sender token
	if err := c.ShouldBindJSON(&body); err != nil { // Transfers request body so that fields match the struct
		log.Println(err.Error())
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Verify the token with the server
	// Function checks if the token is valid and returns the decoded token
	token, err := globals.AuthClient.VerifyIDToken(globals.FirebaseContext, body.IDToken)
	if err != nil {
		log.Error(err.Error())
		c.IndentedJSON(http.StatusForbidden, gin.H{"error": "You are not authorized to report this donation."})
		return
	}
	// Extract sender id from the token
	userUID := token.UID

	// Report the donation using the donationid and the senderid
	err = helpers.ReportDonation(body.DonationID, userUID)
	// If user has already sent a report to this donation, do not continue and send an error to the frontend
	if err != nil {
		log.Error(err.Error())
		if err.Error() == "User has already sent a report" {
			c.IndentedJSON(http.StatusConflict, gin.H{"error": err.Error()})
		} else {
			//if there was some other error, send back a internal server error to the frontend
			log.Println(err.Error())
			c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	c.Status(http.StatusAccepted)
	log.Info("report successful")
}
