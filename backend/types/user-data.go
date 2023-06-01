package types

import (
	"time"

	"cloud.google.com/go/firestore"
)

// UserData represents a user's data.
// It includes display name, email, registration timestamp, admin status, user's posts,
// UID and count of donations made.
type UserData struct {
	DisplayName           string                   `json:"display_name"`
	Email                 string                   `json:"email"`
	RegistrationTimestamp time.Time                `json:"registered_date"` // In UTC
	Admin                 bool                     `json:"admin"`
	Posts                 []*firestore.DocumentRef `json:"posts"`
	UID                   string                   `json:"uid"`
	DonationsMade         int64                    `json:"donations_made"`
}
