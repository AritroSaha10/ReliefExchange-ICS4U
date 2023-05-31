package helpers

// This is a file in the package-"helpers" that contains the GetDonationByID function.
import (
	"fmt"
	"relief_exchange_backend/globals"
	"relief_exchange_backend/types"
	"time"

	log "github.com/sirupsen/logrus"
)

// GetDonationByID retrieves a donation record by its ID from Firestore.
// Parameters:
//   - ctx: the context in which the function is invoked.
//   - client: the Firestore client.
//   - id: the ID of the donation to retrieve.
//
// Return values:
//   - Donation object that corresponds to the provided ID.
//   - error, if any occurred during retrieval.
func GetDonationByID(id string) (types.Donation, error) {
	var donation types.Donation
	doc, err := globals.FirestoreClient.Collection("donations").Doc(id).Get(globals.FirebaseContext) // get a single donation from its id
	if err != nil {
		log.Error(err.Error())
		return donation, err // returns empty donation struct
	}
	err = doc.DataTo(&donation)
	if err != nil {
		log.Error(err.Error())
		return types.Donation{}, err
	}

	// Override some attributes that don't work with DataTo
	data := doc.Data()
	var ok_img, ok_id, ok_time bool

	donation.Image, ok_img = data["img"].(string)
	donation.OwnerId, ok_id = data["owner_id"].(string)
	donation.CreationTimestamp, ok_time = data["creation_timestamp"].(time.Time)
	if !(ok_img && ok_id && ok_time) {
		log.Warn("user data may have not been converted properly")
	}
	// Convert the empty interface types to actual strings
	donation.Reports = make([]string, 0)
	for _, reportRaw := range data["reports"].([]interface{}) {
		donation.Reports = append(donation.Reports, fmt.Sprintf("%+v", reportRaw))
	}

	donation.ID = doc.Ref.ID // ID is stored in the Ref feild, so DataTo, does not store id in the donations object
	log.Info("donation: %v", donation)
	return donation, nil
}
