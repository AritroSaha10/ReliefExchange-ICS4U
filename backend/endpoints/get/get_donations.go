package get

// This file is to modulize the code and contains the GetDonationsList function.

import (
	"net/http"

	"relief_exchange_backend/helpers"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

// getDonationsListEndpoint handles the endpoint to fetch all donations.
// Parameters:
//   - c: the gin context, the request and response http.
//
// It sends a list of all donations in the database to the client.
func GetDonationsList(c *gin.Context) {
	donations, err := helpers.GetAllDonations()
	if err != nil {
		log.Error(err)
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	} else {
		log.Info("Get donations successful.")
		c.IndentedJSON(http.StatusOK, donations)
	}
}
