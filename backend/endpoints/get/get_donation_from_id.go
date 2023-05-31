// Package get provides functions for retrieving information.
package get

// This file is to modulize the code and contains the GetDonationByID function.
import (
	"net/http"
	"relief_exchange_backend/helpers"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

// GetDonationByID handles the endpoint to fetch a donation by id using the getDonationById function
// Parameters:
//   - c: the gin context, the request and response http.
//
// It sends the requested donation to the client.
func GetDonationByID(c *gin.Context) {
	id := c.Param("id")
	donation, err := helpers.GetDonationByID(id)
	if err != nil {
		log.Warn("Donation not found, ID:", id)
		log.Error(err.Error())
		c.IndentedJSON(http.StatusNotFound, gin.H{"error": err.Error()})
	} else {
		log.Info("Get donation by ID successful.")
		c.IndentedJSON(http.StatusOK, donation)
	}
}
