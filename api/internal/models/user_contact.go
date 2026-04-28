package models

import "time"

// UserContact represents a metadata table for active conversations between users.
// Instead of querying millions of messages with DISTINCT, we maintain this table.
type UserContact struct {
	UserID        string    `gorm:"primaryKey;type:uuid;index" json:"user_id"`
	ContactID     string    `gorm:"primaryKey;type:uuid;index" json:"contact_id"`
	LastMessageAt time.Time `json:"last_message_at"`

	Contact *User `gorm:"foreignKey:ContactID;references:ID" json:"contact,omitempty"`
}
