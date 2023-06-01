// @authors Aritro Saha, Joshua Chou
// @cite "Golang: Type Assertion." Better Programming, 2020. [Online].
// Available: https://betterprogramming.pub/golang-type-assertion-d5517d81c366. [Accessed: 23- May- 2023].
// @cite "Package firestore." Pkg.go.dev, 2023. [Online].
// Available: https://pkg.go.dev/cloud.google.com/go/firestore. [Accessed: 22- May- 2023].
// This is a file in the package-"helpers" that contains the BanUser function.
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

	// Check if they're an admin
	isAdmin, err := CheckIfAdmin(userId)
	if err != nil {
		err = fmt.Errorf("err while checking if admin: %w", err)
		log.Error(err.Error())
		return err
	}
	if isAdmin {
		err := fmt.Errorf("cannot ban an admin")
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
