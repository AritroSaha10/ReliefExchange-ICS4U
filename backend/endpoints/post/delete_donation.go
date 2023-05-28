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
	// Extract id of donation from request
	id := c.Param("id")
	// Get the donation firestore document reference
	donationRef := globals.FirestoreClient.Collection("donations").Doc(id)
	// Get the data of the document reference
	donationData, err := donationRef.Get(globals.FirebaseContext)
	if err != nil {
		log.Error(err.Error())
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	// Get the "Authorization" field from the request header
	authHeader := c.GetHeader("Authorization")
	// If the Authorization header is empty, return an error message and a 403 status code
	if authHeader == "" {
		c.JSON(http.StatusForbidden, gin.H{"error": "No authorization header provided"})
		return
	}
	// Split the Authorization header value on the "Bearer" keyword.
	// This is done because the Authorization header usually follows the format "Bearer <token>",
	// where <token> is the actual token value.
	splitToken := strings.Split(authHeader, "Bearer")
	// If the split operation does not result in exactly two parts, return an error message and a 403 status code.
	// This indicates that the Authorization header was not in the expected format. "Bearer <token>"
	if len(splitToken) != 2 {
		c.JSON(http.StatusForbidden, gin.H{"error": "Incorrect format for authorization header"})
		return
	}

	// The TrimPrefix function is used to remove the prefix (Bearer in this case), leaving only the token string.
	tokenString := strings.TrimPrefix(authHeader, "Bearer ")
	// Verify the token using the VerifyIDToken function from the AuthClient.
	// This function checks if the token is valid and returns the decoded token.
	token, err := globals.AuthClient.VerifyIDToken(globals.FirebaseContext, tokenString)
	if err != nil {
		log.Error(err.Error())
		c.IndentedJSON(http.StatusForbidden, gin.H{"error": "You are not authorized to delete this donation."})
		return
	}

	// Extract the user id from the token
	userUID := token.UID
	// Only allow donation owner, or admins to delete this donation
	// If sender id (userUID) does not match the id of the donation owner, or the sender id, is not an admin, then they are not authorized to delete the donation
	isAdmin, _ := helpers.CheckIfAdmin(userUID)
	if donationData.Data()["owner_id"].(string) != userUID && (!isAdmin) {
		c.JSON(http.StatusForbidden, gin.H{"error": "You are not authorized to delete this donation."})
		return
	}

	// Delete the donation using the firestore reference
	_, err = donationRef.Delete(globals.FirebaseContext) // only need the err return value
	if err != nil {
		log.Error(err.Error())
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	// If the deletion was successful, return a 200 OK status and a success message.
	c.JSON(http.StatusOK, gin.H{"message": "Donation deleted successfully"})
}
