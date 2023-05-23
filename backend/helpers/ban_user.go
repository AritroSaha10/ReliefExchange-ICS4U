package helpers

import (
	"fmt"
	"relief_exchange_backend/globals"

	"cloud.google.com/go/firestore"
	log "github.com/sirupsen/logrus"
)

// BanUser bans a user by removing their records from Firestore and flagging their UID.
// Parameters:
//   - ctx: the context in which the function is invoked.
//   - client: the Firestore client.
//   - userId: the ID of the user to ban.
//
// Return values:
//   - error, if any occurred during the operation.
func BanUser(userId string) error {
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

	// Delete user data
	if _, err := userDataRef.Delete(globals.FirebaseContext); err != nil {
		err = fmt.Errorf("failed deleting user data: %w", err)
		log.Error(err.Error())
		return err
	}

	// Add them to the banned list
	banDocRef := globals.FirestoreClient.Doc("config/bans")
	var banDocSnapshot *firestore.DocumentSnapshot
	if banDocSnapshot, err = banDocRef.Get(globals.FirebaseContext); err != nil {
		err = fmt.Errorf("failed getting ban list: %w", err)
		log.Error(err.Error())
		return err
	}

	// Get ban list
	var banList []string
	if banList, ok = banDocSnapshot.Data()["users"].([]string); !ok {
		err = fmt.Errorf("could not convert banned users list to []string")
		log.Error(err.Error())
		return err
	}
	banList = append(banList, userId)

	// Update the document with new banned user
	banDocRef.Update(globals.FirebaseContext, []firestore.Update{
		{
			Path:  "users",
			Value: banList,
		},
	})

	return nil
}
