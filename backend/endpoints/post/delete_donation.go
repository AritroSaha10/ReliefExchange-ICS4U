package post

import (
	"net/http"
	"relief_exchange_backend/globals"
	"relief_exchange_backend/helpers"
	"strings"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

// DeleteDonation handles the endpoint to delete a donation by id.
// Parameters:
//   - c: the gin context, the request and response http.
//
// It requires an Authorization header with a bearer token, verifies the token,
// checks if the user is authorized to delete the donation, then deletes it.
func DeleteDonation(c *gin.Context) {
	id := c.Param("id")

	donationRef := globals.FirestoreClient.Collection("donations").Doc(id)
	donationData, err := donationRef.Get(globals.FirebaseContext)
	if err != nil {
		log.Error(err.Error())
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	authHeader := c.GetHeader("Authorization")

	if authHeader == "" {
		c.JSON(http.StatusForbidden, gin.H{"error": "No authorization header provided"})
		return
	}

	splitToken := strings.Split(authHeader, "Bearer")

	if len(splitToken) != 2 {
		c.JSON(http.StatusForbidden, gin.H{"error": "Incorrect format for authorization header"})
		return
	}

	tokenString := strings.TrimPrefix(authHeader, "Bearer ")

	token, err := globals.AuthClient.VerifyIDToken(globals.FirebaseContext, tokenString)
	if err != nil {
		log.Error(err.Error())
		c.IndentedJSON(http.StatusForbidden, gin.H{"error": "You are not authorized to delete this donation."})
		return
	}

	userUID := token.UID

	isAdmin, _ := helpers.CheckIfAdmin(userUID)
	if donationData.Data()["owner_id"].(string) != userUID && (!isAdmin) {
		c.JSON(http.StatusForbidden, gin.H{"error": "You are not authorized to delete this donation."})
		return
	}

	_, err = donationRef.Delete(globals.FirebaseContext) // only need the err return value
	if err != nil {
		log.Error(err.Error())
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Donation deleted successfully"})
}
