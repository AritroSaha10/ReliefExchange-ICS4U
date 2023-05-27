package post

import (
	"net/http"
	"relief_exchange_backend/globals"
	"relief_exchange_backend/helpers"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

// BanUser handles the endpoint to ban a user.
// Parameters:
//   - c: the gin context, the request and response http.
//
// It accepts a user's id token and the id of the user to be banned, verifies the token,
// checks if the user performing the ban is an admin, then bans the user if authorized using the banUser function and the checkIfAdmin function
func DeleteProfile(c *gin.Context) {
	var body struct {
		ProfToDelete string `json:"profToDelete"`
		Token        string `json:"token"`
	}

	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// get sending user token
	token, err := globals.AuthClient.VerifyIDToken(globals.FirebaseContext, body.Token)
	if err != nil {
		log.Error(err.Error())
		c.IndentedJSON(http.StatusForbidden, gin.H{"error": "you are not authorized to delete this user"})
		return
	}

	// get uuid of user to ban
	uuidToDelete := body.ProfToDelete
	log.Info(uuidToDelete)

	if uuidToDelete == token.UID {
		// if sending user is an admin delete all the donations of the user to ban including their data and account
		err = helpers.DeleteProfile(uuidToDelete)
		if err != nil {
			log.Error(err.Error())
			c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": "There was an error processing the deletion"})
			return
		}

		c.IndentedJSON(http.StatusOK, gin.H{"status": "User deleted successfully"})
	} else {
		log.Error(err.Error())
		c.IndentedJSON(http.StatusForbidden, gin.H{"error": "You are not authorized to delete this user"})
	}
}
