package helpers

import (
	"fmt"
	"relief_exchange_backend/globals"

	log "github.com/sirupsen/logrus"
)

// CheckIfAdmin checks if a user has admin privileges.
// Parameters:
//   - senderId: the ID of the user to check.
//
// Return values:
//   - true if the user has admin privileges, false otherwise.
//   - error, if any occurred during the check.
func CheckIfAdmin(senderId string) (bool, error) {
	// Get the user document
	doc, err := globals.FirestoreClient.Doc("users/" + senderId).Get(globals.FirebaseContext)

	if err != nil {
		log.Error(err.Error())
		return false, err
	}

	// Get the data and access the 'isAdmin' field
	data := doc.Data()

	isAdmin, ok := data["admin"].(bool)
	if !ok {
		err = fmt.Errorf("failed getting admin field from user doc: %w", err)
		log.Error(err.Error())
		return false, err
	}

	log.Info("isAdmin: %v", isAdmin)
	return isAdmin, nil
}
