package models

import (
	"time"

	"gorm.io/gorm"
)

type Message struct {
	ID      string `json:"id" gorm:"primaryKey"`
	Content string `json:"content"`

	FromID string `json:"from_id" gorm:"column:from_id;not null;index:idx_conversation,priority:1"`
	ToID   string `json:"to_id" gorm:"column:to_id;not null;index:idx_conversation,priority:2;index:idx_unseen,priority:1,where:seen_at IS NULL"`

	DeliveredAt *time.Time `json:"delivered_at,omitempty"`
	SeenAt      *time.Time `json:"seen_at,omitempty"`

	CreatedAt time.Time      `json:"created_at" gorm:"index:idx_conversation,priority:3,sort:desc"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"deleted_at,omitempty" gorm:"index"`
}

type GroupMessage struct {
	ID      string `json:"id" gorm:"primaryKey"`
	Content string `json:"content"`

	FromID string `json:"from_id" gorm:"column:from_id;index;not null"`

	ToID string `json:"to_id" gorm:"column:to_id;index;not null"`

	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at" gorm:"index"`
}

type MessageQueue struct {
	ID        string     `json:"id" gorm:"primaryKey"`
	MessageID string     `json:"message_id" gorm:"index"`
	GroupID   string     `json:"group_id" gorm:"index"`
	UserID    string     `json:"user_id" gorm:"index"`
	CreatedAt time.Time  `json:"created_at" gorm:"index"`
	UpdatedAt time.Time  `json:"updated_at" gorm:"index"`
	DeletedAt *time.Time `json:"deleted_at" gorm:"index"`
}
