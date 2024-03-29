package get

// This file is to modulize the code and contains the GetIfBanned function.
// @author Aritro Saha
import (
	"net/http"
	"relief_exchange_backend/helpers"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

// banUserEndpoint handles the endpoint to check if a user is banned.
// Parameters:
//   - c: the gin context, the request and response http.
//
// It accepts a user's id token and the id of the user to be checked, and checks
// if they have been banned on the platform.
func GetIfBanned(c *gin.Context) {
	userUID := c.Query("uid")

	// Get the result from the helper function
	isBanned, err := helpers.CheckIfBanned(userUID)
	if err != nil {
		log.Error(err.Error())
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	// Return result to user
	c.IndentedJSON(http.StatusOK, gin.H{"banned": isBanned})
}
