package post

import (
	"net/http"
	"relief_exchange_backend/globals"
	"relief_exchange_backend/helpers"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

// AddUser handles the endpoint to add a new user using the addUser function.
// Parameters:
//   - c: the gin context, the request and response http.
//
// It accepts a user's id token, verifies the token, and then adds the user to the database.
func AddUser(c *gin.Context) {
	var body struct {
		IDToken string `json:"token"`
	}
	// Attempt to bind the JSON body of the request to the struct
	if err := c.ShouldBindJSON(&body); err != nil {
		log.Error(err.Error())
		// stores request body info into the body varible, so that it matches feild in struct in json format
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

	// Attempt to add the user to the database using the AddUser functions
	err = helpers.AddUser(userUID)
	// Add to the firestore databse
	if err != nil {
		log.Error(err.Error())
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	} else {
		log.Info("user added successfully")
		c.IndentedJSON(http.StatusCreated, gin.H{"message": "User added successfully"})
	}
}
