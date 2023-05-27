package helpers

import (
	"fmt"
	"relief_exchange_backend/globals"

	"cloud.google.com/go/firestore"
	log "github.com/sirupsen/logrus"
)

// BanUser bans a user by removing their records from Firestore and flagging their UID.
// Parameters:
//   - userId: the ID of the user to ban.
//
// Return values:
//   - error, if any occurred during the operation.
func DeleteProfile(userId string) error {

	// Get user data document reference
	userDataRef := globals.FirestoreClient.Doc("users/" + userId)

	// Get user data
	userDataDoc, err := userDataRef.Get(globals.FirebaseContext)
	if err != nil {
		err = fmt.Errorf("failed getting user data: %w", err)
		log.Error(err.Error())
		return err
	}

	// Extract posts field from user data
	// Note: .([]interface{}) is a type assertion that checks if the value returned is a slice of interfaces,
	// The attributes in doc.Data() are assumed to be either of an unknown type or a slice of unknown type.
	// If it were a single unknown type, we could directly convert it to the desired type.
	// However, we can't directly convert a slice of unknown type to a slice of a specific type.
	// Therefore, we have convert every element in the slice to the desired type,
	// and then add it to a new slice of the desired type.
	rawPosts, ok := userDataDoc.Data()["posts"].([]interface{})
	if !ok {
		err = fmt.Errorf("failed extracting posts field from user data")
		log.Error(err.Error())
		return err
	}

	// Convert each raw post to a *firestore.DocumentRef
	for _, rawPost := range rawPosts {
		postRef, ok := rawPost.(*firestore.DocumentRef)
		if !ok {
			log.Warn("failed converting raw post to *firestore.DocumentRef")
			continue
		}
		if _, err := postRef.Delete(globals.FirebaseContext); err != nil {
			log.Warn("failed deleting post: %w", err)
			continue
		}
	}

	return nil
}
