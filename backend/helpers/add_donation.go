// @cite "How to insert document id as a field value Firestore with Go." Stack Overflow, 2020. [Online].
// Available: https://stackoverflow.com/questions/61207401/how-to-insert-document-id-as-a-field-value-firestore-with-go. [Accessed: 27- April- 2023].
package helpers

import (
	"fmt"
	"relief_exchange_backend/globals"
	"relief_exchange_backend/types"

	"cloud.google.com/go/firestore"
	log "github.com/sirupsen/logrus"
)

// AddDonation adds a new donation record to Firestore.
// Parameters:
//   - ctx: the context in which the function is invoked.
//   - client: the Firestore client.
//   - donation: the Donation object to add.
//   - userId: the ID of the user making the donation.
//
// Return values:
//   - ID of the new donation record.
//   - error, if any occurred during the operation.
func AddDonation(donation types.Donation, userId string) (string, error) {
	// Check if they're already banned
	banned, err := CheckIfBanned(userId)
	if err != nil {
		err = fmt.Errorf("err while checking if banned: %w", err)
		log.Error(err.Error())
		return "", err
	}
	if banned {
		err := fmt.Errorf("user is already banned")
		log.Error(err.Error())
		return "", err
	}

	docRef, _, err := globals.FirestoreClient.Collection("donations").Add(globals.FirebaseContext, map[string]interface{}{
		"title":              donation.Title,
		"description":        donation.Description,
		"location":           donation.Location,
		"img":                donation.Image,
		"owner_id":           userId,
		"creation_timestamp": donation.CreationTimestamp,
		"tags":               donation.Tags,
		"reports":            make([]string, 0),
	})
	if err != nil {
		err = fmt.Errorf("error while adding donation: %w", err)
		log.Error(err.Error())
		return "", err
	}

	// Get current posts and append new post
	// Get the user's document
	userDoc, err := globals.FirestoreClient.Doc("users/" + userId).Get(globals.FirebaseContext)
	if err != nil {
		err = fmt.Errorf("error while getting user document (addDonation): %w", err)
		log.Error(err.Error())
		return "", err
	}

	// Extract posts array
	data := userDoc.Data()
	rawPosts := data["posts"].([]interface{})

	// Iterate over the array field to find document references
	var posts []*firestore.DocumentRef
	for _, value := range rawPosts {
		posts = append(posts, value.(*firestore.DocumentRef))
	}
	// Append the new donation reference
	posts = append(posts, docRef)

	// Update the user document
	_, err = globals.FirestoreClient.Doc("users/"+userId).Set(globals.FirebaseContext, map[string]interface{}{
		"posts":          posts,
		"donations_made": userDoc.Data()["donations_made"].(int64),
	}, firestore.MergeAll) // mergeall ensures that only the posts feild is changed
	if err != nil {
		err = fmt.Errorf("error while updating user document (addDonation): %w", err)
		log.Error(err.Error())
		return "", err
	}

	log.Info("ID of new donation: %v", docRef.ID)
	return docRef.ID, nil
}
