package helpers

import (
	"relief_exchange_backend/globals"
	"relief_exchange_backend/types"
	"time"

	log "github.com/sirupsen/logrus"
	"google.golang.org/api/iterator"
)

// getAllDonations retrieves all donation records from Firestore.
// Parameters:
//   - ctx: the context in which the function is invoked.
//   - client: the Firestore client.
//
// Return values:
//   - Slice of all Donation objects retrieved.
//   - error, if any occurred during retrieval.
func GetAllDonations() ([]types.Donation, error) {
	var donations []types.Donation
	iter := globals.FirestoreClient.Collection("donations").Documents(globals.FirebaseContext) //.Documents(ctx) returns a iterator
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			log.Error(err.Error())
			return nil, err // no data was retrieved-nil, but there was an error -err
		}
		var donation types.Donation
		err = doc.DataTo(&donation)

		// Override some attributes that don't work with DataTo
		donation.Image = doc.Data()["img"].(string)
		donation.OwnerId = doc.Data()["owner_id"].(string)
		donation.CreationTimestamp = doc.Data()["creation_timestamp"].(time.Time)

		if err != nil {
			log.Error(err.Error())
			return nil, err
		}
		donation.ID = doc.Ref.ID // sets donation struct id to the one in the firebase
		donations = append(donations, donation)
	}

	log.Info("donations:%v", donations)

	return donations, nil // nil-data was retrived without any errors
}
