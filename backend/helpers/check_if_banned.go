package helpers

import (
	"context"
	"fmt"
	"relief_exchange_backend/globals"

	"cloud.google.com/go/firestore"
	log "github.com/sirupsen/logrus"
	"golang.org/x/exp/slices"
)

// checkIfBanned checks whether a user is banned by looking at the config docs in Firestore.
// Parameters:
//   - ctx: the context in which the function is invoked.
//   - client: the Firestore client.
//   - userId: the ID of the user to check.
//
// Return values:
//   - bool, if they are banned or not
//   - error, if any occurred during the operation.
func CheckIfBanned(ctx context.Context, userId string) (bool, error) {
	// Add them to the banned list
	banDocRef := globals.FirestoreClient.Doc("config/bans")
	var banDocSnapshot *firestore.DocumentSnapshot
	var err error
	if banDocSnapshot, err = banDocRef.Get(ctx); err != nil {
		log.Error(err.Error())
		return false, fmt.Errorf("failed getting ban list: %w", err)
	}

	// Get ban list
	var banListRaw []interface{}
	var ok bool
	if banListRaw, ok = banDocSnapshot.Data()["users"].([]interface{}); !ok {
		err = fmt.Errorf("could not convert banned users list to []string")
		log.Error(err.Error())
		return false, err
	}

	var banList []string
	for _, rawBannedUID := range banListRaw {
		if bannedUID, ok := rawBannedUID.(string); !ok {
			log.Warn("could not convert a UID in banned list to string")
			continue
		} else {
			banList = append(banList, bannedUID)
		}
	}

	// Return whether uid in list
	return slices.Contains(banList, userId), nil
}
