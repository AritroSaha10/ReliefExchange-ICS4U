package helpers

import (
	"fmt"
	"relief_exchange_backend/globals"

	"cloud.google.com/go/firestore"
	log "github.com/sirupsen/logrus"
)

// ReportDonation adds a report to a specific donation record.
// Parameters:
//   - donationID: the ID of the donation to report.
//   - userUID: the UID of the user making the report.
//
// Return values:
//   - error, if any occurred during the operation.
func ReportDonation(donationID string, userUID string) error {
	doc, err := globals.FirestoreClient.Collection("donations").Doc(donationID).Get(globals.FirebaseContext) // Get the donation's data
	if err != nil {
		log.Error(err.Error())
		return err
	}

	// Check whether they have already reported this donation
	currentReports, ok := doc.Data()["reports"].([]string)
	if !ok {
		// Convert the empty interface types to actual strings
		currentReports = make([]string, 0)
		for _, reportRaw := range doc.Data()["reports"].([]interface{}) {
			currentReports = append(currentReports, fmt.Sprintf("%+v", reportRaw))
		}
	}

	// Check whether they've already made a report
	for _, report := range currentReports {
		if report == userUID {
			err := fmt.Errorf("user has already sent a report")
			log.Error(err)
			return err
		}
	}

	// Add their UID to the donation, and update the doc
	newReports := append(currentReports, userUID)
	_, err = globals.FirestoreClient.Collection("donations").Doc(donationID).Update(globals.FirebaseContext, []firestore.Update{
		{
			Path:  "reports",
			Value: newReports,
		},
	})
	if err != nil {
		err = fmt.Errorf("failed adding report to donation doc: %w", err)
		log.Error(err.Error())
		return err
	}

	return nil
}