package helpers

import (
	"context"
	"relief_exchange_backend/globals"
	"relief_exchange_backend/types"
	"time"

	log "github.com/sirupsen/logrus"
)

// getUserDataByID retrieves user data by the user's ID from Firestore.
// Parameters:
//   - ctx: the context in which the function is invoked.
//   - client: the Firestore client.
//   - id: the ID of the user to retrieve.
//
// Return values:
//   - UserData object that corresponds to the provided ID.
//   - error, if any occurred during retrieval.
func GetUserDataByID(ctx context.Context, id string) (types.UserData, error) {
	var userData types.UserData
	doc, err := globals.FirestoreClient.Collection("users").Doc(id).Get(ctx) // Get a single user from its id
	if err != nil {
		log.Error(err.Error())
		return userData, err // returns empty user struct
	}
	err = doc.DataTo(&userData)
	if err != nil {
		log.Error(err.Error())
		return types.UserData{}, err
	}

	// Set values that aren't set in the DataTo function
	var ok1, ok2, ok3 bool
	userData.DisplayName, ok1 = doc.Data()["display_name"].(string)
	userData.RegistrationTimestamp, ok2 = doc.Data()["registered_date"].(time.Time)
	userData.DonationsMade, ok3 = doc.Data()["donations_made"].(int64)
	if !(ok1 && ok2 && ok3) {
		log.Warn("user data may have not been converted properly")
	}

	userData.UID = doc.Ref.ID // ID is stored in the Ref feild, so DataTo, does not store id in the user data object
	log.Info("userData: %v", userData)
	return userData, nil
}
