package helpers

// This is a file in the package-"helpers" that contains the GerUserDataByID function.
import (
	"fmt"
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
func GetUserDataByID(id string) (types.UserData, error) {
	// Check if they're already banned
	banned, err := CheckIfBanned(id)
	if err != nil {
		err = fmt.Errorf("err while checking if banned: %w", err)
		log.Error(err.Error())
		return types.UserData{}, err
	}
	if banned {
		err := fmt.Errorf("user is banned")
		log.Error(err.Error())
		return types.UserData{}, err
	}

	var userData types.UserData
	doc, err := globals.FirestoreClient.Collection("users").Doc(id).Get(globals.FirebaseContext) // Get a single user from its id
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
	var ok_name, ok_date, ok_donations_made bool
	userData.DisplayName, ok_name = doc.Data()["display_name"].(string)
	userData.RegistrationTimestamp, ok_date = doc.Data()["registered_date"].(time.Time)
	userData.DonationsMade, ok_donations_made = doc.Data()["donations_made"].(int64)
	if !(ok_name && ok_date && ok_donations_made) {
		log.Warn("user data may have not been converted properly")
	}

	userData.UID = doc.Ref.ID // ID is stored in the Ref feild, so DataTo, does not store id in the user data object
	log.Info("userData: %v", userData)
	return userData, nil
}
