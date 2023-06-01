package helpers

// @author Aritro Saha
// This is a file in the package-"helpers" that contains the EditDonation function.
import (
	"fmt"
	"relief_exchange_backend/globals"
	"relief_exchange_backend/types"

	log "github.com/sirupsen/logrus"
)

// EditDonation edits an existing donation record on Firestore.
// Parameters:
//   - newDonation: the new Donation data.
//   - currId: the ID of the current donation
//
// Return values:
//   - error, if any occurred during the operation.
func EditDonation(newDonation types.Donation, currId string) error {
	// Get a reference to the current Donation doc
	docRef := globals.FirestoreClient.Collection("donations").Doc(currId)
	oldData, err := docRef.Get(globals.FirebaseContext)
	if err != nil {
		err = fmt.Errorf("err while getting current donation ref: %w", err)
		log.Error(err.Error())
		return err
	}

	// Change current donation data to new data
	_, err = docRef.Set(globals.FirebaseContext, map[string]interface{}{
		"title":              newDonation.Title,
		"description":        newDonation.Description,
		"location":           newDonation.Location,
		"img":                oldData.Data()["img"].(string), // Don't allow editing photos
		"owner_id":           oldData.Data()["owner_id"].(string),
		"creation_timestamp": newDonation.CreationTimestamp,
		"tags":               newDonation.Tags,
		"reports":            make([]string, 0),
	})
	if err != nil {
		err = fmt.Errorf("error while updating donation: %w", err)
		log.Error(err.Error())
		return err
	}

	return nil
}
