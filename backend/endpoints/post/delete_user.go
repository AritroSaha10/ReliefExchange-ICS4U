/*
 * File: delete_user.go
 * -------------
 * This module handles the delete user endpoint in the server.
 * It takes a gin context as a parameter, binds the request body to a struct,
 * extracts the token from it, and verifies the token.
 * If the token is valid, it calls the DeleteUser helper function to delete the user with the uid extracted from the token.
 */
// @author Joshua Chou
// @cite "Validating Google Sign In ID Token in Go." Stack Overflow, 2016. [Online].
// Available: https://stackoverflow.com/questions/36716117/validating-google-sign-in-id-token-in-go. [Accessed: 27- May- 2023].
package post

import (
	"net/http"
	"relief_exchange_backend/globals"
	"relief_exchange_backend/helpers"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

// DeleteUser handles the endpoint to delete all of a user's data
// using the deleteUser function.
// Parameters:
//   - c: the gin context, the request and response http.
//
// It accepts a user's id token, verifies the token, and then deletes all of their data.
func DeleteUser(c *gin.Context) {
	var body struct {
		IDToken string `json:"token"`
	}
	// Attempt to bind the JSON body of the request to the struct
	if err := c.ShouldBindJSON(&body); err != nil {
		log.Error(err.Error())
		// Stores request body info into the body varible, so that it matches field in struct in json format
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()}) // if user not signed in, then will send error
		return
	}

	// Attempt to verify the ID token
	// Token is provided for user to verify themselves with the server
	// After it is decoded, we have access to all fields
	token, err := globals.AuthClient.VerifyIDToken(globals.FirebaseContext, body.IDToken)
	if err != nil {
		log.Error(err.Error())
		c.IndentedJSON(http.StatusForbidden, gin.H{"error": "You are not authorized to create this user"})
		return
	}
	// Extract the user's UID from the token
	userUID := token.UID

	// Attempt to delete the user using the DeleteUser function
	err = helpers.DeleteUser(userUID)
	if err != nil {
		log.Error(err.Error())
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	} else {
		log.Info("user added successfully")
		c.IndentedJSON(http.StatusCreated, gin.H{"message": "User deleted successfully"})
	}
}
