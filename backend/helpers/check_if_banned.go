// @cite "How to insert a reference type field on Firestore with Golang." Stack Overflow, 2021. [Online].
// Available: https://stackoverflow.com/questions/69797221/how-to-insert-a-reference-type-field-on-firestore-with-golang. [Accessed: 18- May- 2023].
// This is a file in the package-"helpers" that contains the CheckIfBanned function.
package helpers

import (
	"fmt"
	"relief_exchange_backend/globals"

	"cloud.google.com/go/firestore"
	log "github.com/sirupsen/logrus"
	"golang.org/x/exp/slices"
)

// checkIfBanned checks whether a user is banned by looking at the config docs in Firestore.
// Parameters:
//   - client: the Firestore client.
//   - userId: the ID of the user to check.
//
// Return values:
//   - bool, if they are banned or not
//   - error, if any occurred during the operation.
func CheckIfBanned(userId string) (bool, error) {
	// Add them to the banned list
	banDocRef := globals.FirestoreClient.Doc("config/bans")
	var banDocSnapshot *firestore.DocumentSnapshot
	var err error
	if banDocSnapshot, err = banDocRef.Get(globals.FirebaseContext); err != nil {
		log.Error(err.Error())
		return false, fmt.Errorf("failed getting ban list: %w", err)
	}

	// Get ban list

	// Get raw ban list directly from Firestore
	var banListRaw []interface{}
	var ok bool
	if banListRaw, ok = banDocSnapshot.Data()["users"].([]interface{}); !ok {
		err = fmt.Errorf("ban users list is not of type []interface{}")
		log.Error(err.Error())
		return false, err
	}

	// Convert raw ban list of type []interface{} to type []string
	var banList []string
	for _, rawBannedUID := range banListRaw {
		if bannedUID, ok := rawBannedUID.(string); !ok {
			log.Warn("UID in banned list is not a string")
			continue
		} else {
			banList = append(banList, bannedUID)
		}
	}

	// Return whether uid in list
	return slices.Contains(banList, userId), nil
}
