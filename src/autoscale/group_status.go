package autoscale

import "time"

// GroupStatus is a log of a scaling event for a group.
type GroupStatus struct {
	GroupID   string    `json:"groupID" db:"group_id"`
	Delta     int       `json:"delta" db:"delta"`
	Total     int       `json:"total" db:"total"`
	CreatedAt time.Time `json:"createdAt" db:"created_at"`
}
