/*
 * File: edit_donation.go
 * -------------
 * This module handles the edit donation endpoint in the server.
 * It takes a gin context as a parameter, binds the request body to a struct,
 * extracts the token, donation data and donation id from it, and verifies the token.
 * If the token is valid and the owner id from the token matches the owner id in the donation data or the user is an admin,
 * it calls the EditDonation helper function to edit the donation with the given donation id and the new donation data.
 */

package post

import (
	"fmt"
	"net/http"
	"relief_exchange_backend/globals"
	"relief_exchange_backend/helpers"
	"relief_exchange_backend/types"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

// EditDonation handles the endpoint to edit an existing donation.
// Parameters:
//   - c: the gin context, the request and response http.
//
// It accepts an existing donation UID, new donation data, and a user's id token,
// verifies the token, and then uses the EditDonation helper to edit the donation
// in the database.
func EditDonation(c *gin.Context) {
	var body struct {
		ExistingDonationID string         `json:"id"`
		NewDonationData    types.Donation `json:"data"`
		IDToken            string         `json:"token"`
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

	// Extract user data from token verification
	userUID := token.UID

	// Get donation data to extract creator's UID
	existingDonation, err := helpers.GetDonationByID(body.ExistingDonationID)
	if err != nil {
		err = fmt.Errorf("err while getting existing donation: %w", err)
		log.Error(err.Error())
		c.Status(http.StatusInternalServerError)
		return
	}

	// Check whether user is an admin
	isAdmin, err := helpers.CheckIfAdmin(userUID)
	if err != nil {
		err = fmt.Errorf("err while checking if admin: %w", err)
		log.Error(err.Error())
		c.Status(http.StatusInternalServerError)
		return
	}

	// Only allow the original creator or an admin to edit posts
	if !(isAdmin || existingDonation.OwnerId == userUID) {
		err = fmt.Errorf("user cannot edit donation, is not the original author or an admin")
		log.Error(err.Error())
		c.IndentedJSON(http.StatusForbidden, gin.H{"error": "request user is not author or admin"})
	}

	// Check if they're already banned
	banned, err := helpers.CheckIfBanned(userUID)
	if err != nil {
		err = fmt.Errorf("err while checking if banned: %w", err)
		log.Error(err.Error())
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if banned {
		err := fmt.Errorf("user is banned")
		log.Error(err.Error())
		c.IndentedJSON(http.StatusForbidden, gin.H{"error": "you were banned from the platform"})
		return
	}

	// Use EditDonation function to edit the donation, passing in the donationData and the existing UID
	err = helpers.EditDonation(body.NewDonationData, body.ExistingDonationID)

	// If there's an error adding the donation, send back err msg to frontend,
	// otherwise send back docId for the frontend to use
	if err != nil {
		log.Error(err.Error())
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	} else {
		log.Info("Editing donation successful.")
		c.Status(http.StatusOK)
	}
}
