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

	if err := c.ShouldBindJSON(&body); err != nil {
		log.Error(err.Error())
		// stores request body info into the body varible, so that it matches feild in struct in json format
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()}) // if user not signed in, then will send error
		return
	}

	token, err := globals.AuthClient.VerifyIDToken(globals.FirebaseContext, body.IDToken) // token is for user to verify with the server, after it is decoded, we have access to all feilds
	if err != nil {
		log.Error(err.Error())
		c.IndentedJSON(http.StatusForbidden, gin.H{"error": "You are not authorized to create this user"})
		return
	}

	userUID := token.UID

	err = helpers.AddUser(userUID) // create new donation object from struct
	// add to the firestore databse
	if err != nil {
		log.Error(err.Error())
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	} else {
		log.Info("user added successfully")
		c.IndentedJSON(http.StatusCreated, gin.H{"message": "User added successfully"})
	}
}
