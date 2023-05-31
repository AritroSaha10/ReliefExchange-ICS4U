package get

// This file is to modulize the code and contains the GetIfAdmin function.

import (
	"net/http"
	"relief_exchange_backend/helpers"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

// GetIfAdmin handles the endpoint to check if a user is an admin.
// Parameters:
//   - c: the gin context, the request and response http.
//
// It accepts a user's UID, and checks
// if they are an admin or not.
func GetIfAdmin(c *gin.Context) {
	userUID := c.Query("uid")

	// Get the result from the helper function
	isAdmin, err := helpers.CheckIfAdmin(userUID)
	if err != nil {
		log.Error(err.Error())
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	// Return result to user
	c.IndentedJSON(http.StatusOK, gin.H{"admin": isAdmin})
}
