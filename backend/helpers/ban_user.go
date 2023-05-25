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
func BanUser(userId string) error {
	// Check if they're already banned
	banned, err := CheckIfBanned(userId)
	if err != nil {
		err = fmt.Errorf("err while checking if banned: %w", err)
		log.Error(err.Error())
		return err
	}
	if banned {
		err := fmt.Errorf("user is already banned")
		log.Error(err.Error())
		return err
	}

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
	if rawBanList, ok := banDocSnapshot.Data()["users"].([]interface{}); !ok {
		err = fmt.Errorf("could not convert banned users list to []interface{}")
		log.Error(err.Error())
		return err
	} else {
		for _, uid := range rawBanList {
			banList = append(banList, uid.(string))
		}
	}

	// Add user to ban list
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
