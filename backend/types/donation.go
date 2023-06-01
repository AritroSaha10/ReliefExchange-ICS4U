package types

import (
	"time"
)

// Donation represents a donation item.
// It includes information about the item like title, description, location, image,
// creation timestamp, owner's id, tags, and reports.
type Donation struct {
	ID                string    `json:"id"`
	Title             string    `json:"title"`
	Description       string    `json:"description"`
	Location          string    `json:"location"`
	Image             string    `json:"img"`
	CreationTimestamp time.Time `json:"creation_timestamp"` // In UTC
	OwnerId           string    `json:"owner_id"`
	Tags              []string  `json:"tags"`
	Reports           []string  `json:"reports"` // Includes the UIDs of every person who reported it
}
