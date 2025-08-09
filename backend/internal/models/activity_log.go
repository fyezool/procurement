package models

import "time"

// ActivityLog represents a recorded action in the system.
type ActivityLog struct {
	ID          int       `json:"id"`
	UserID      *int      `json:"user_id,omitempty"` // Use pointer for nullable foreign key
	Action      string    `json:"action"`
	TargetType  *string   `json:"target_type,omitempty"` // e.g., "USER", "VENDOR"
	TargetID    *int      `json:"target_id,omitempty"`   // e.g., the ID of the affected user or vendor
	Status      string    `json:"status"`                // e.g., "SUCCESS", "FAILED"
	Details     *string   `json:"details,omitempty"`     // e.g., error message on failure
	CreatedAt   time.Time `json:"created_at"`
}
