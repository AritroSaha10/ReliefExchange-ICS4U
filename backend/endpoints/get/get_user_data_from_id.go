package get

// This file is to modulize the code and contains the GetUserDataByID function.

import (
	"net/http"
	"relief_exchange_backend/helpers"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

// GetUserDataByID handles the endpoint to fetch a user's data by id using the helpers.GetUserDataByID Function
// Parameters:
//   - c: the gin context, the request and response http.
//
// It sends the requested user's data to the client.
func GetUserDataByID(c *gin.Context) {
	id := c.Param("id")
	userData, err := helpers.GetUserDataByID(id)
	if err != nil {
		log.Warn("User data not found, ID:", id)
		log.Error(err.Error())
		c.IndentedJSON(http.StatusNotFound, gin.H{"error": err.Error()})
	} else {
		log.Info("Get user data by ID successful.")
		c.IndentedJSON(http.StatusOK, userData)
	}
}
