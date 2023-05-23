package helpers

import (
	"fmt"
	"relief_exchange_backend/globals"
	"time"

	"cloud.google.com/go/firestore"
	log "github.com/sirupsen/logrus"
)

// addUser adds a new user to Firestore.
// Parameters:
//   - ctx: the context in which the function is invoked.
//   - client: the Firestore client.
//   - userId: the ID of the user to add.
//
// Return values:
//   - error, if any occurred during the operation.
func AddUser(userId string) error {
	userData, err := globals.AuthClient.GetUser(globals.FirebaseContext, userId)
	if err != nil {
		err = fmt.Errorf("failed getting user data from auth server: %w", err)
		log.Error(err.Error())
		return err
	}

	_, err = globals.FirestoreClient.Doc("users/"+userId).Create(globals.FirebaseContext, map[string]interface{}{
		"display_name":    userData.DisplayName,
		"email":           userData.Email,
		"admin":           false,
		"posts":           []firestore.DocumentRef{},
		"uid":             userId,
		"donations_made":  0,
		"registered_date": time.Unix(userData.UserMetadata.CreationTimestamp/1000, 0),
	})
	if err != nil {
		err = fmt.Errorf("failed creating user data doc: %w", err)
		log.Error(err.Error())
		return err
	}
	return nil
}
