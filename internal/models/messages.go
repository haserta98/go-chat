package models

import (
	"time"

	"gorm.io/gorm"
)

type Message struct {
	ID      string `json:"id" gorm:"primaryKey"`
	Content string `json:"content"`

	FromID   string `json:"from_id" gorm:"column:from_user;not null;index:idx_conversation,priority:1"`
	FromUser User   `json:"from_user" gorm:"foreignKey:FromID;references:ID"`

	ToUserID string `json:"to_user_id" gorm:"not null;index:idx_conversation,priority:2;index:idx_unseen,priority:1,where:seen_at IS NULL"`
	ToUser   *User  `json:"to_user,omitempty" gorm:"foreignKey:ToUserID;references:ID"`

	DeliveredAt *time.Time `json:"delivered_at,omitempty"`
	SeenAt      *time.Time `json:"seen_at,omitempty"`

	CreatedAt time.Time      `json:"created_at" gorm:"index:idx_conversation,priority:3,sort:desc"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"deleted_at,omitempty" gorm:"index"`
}

type GroupMessage struct {
	ID      string `json:"id" gorm:"primaryKey"`
	Content string `json:"content"`

	FromID   string `json:"from_id" gorm:"column:from_user;index;not null"`
	FromUser User   `json:"from_user" gorm:"foreignKey:FromID;references:ID"`

	ToGroupID string `json:"to_group_id,omitempty" gorm:"index"`
	ToGroup   *Group `json:"to_group,omitempty" gorm:"foreignKey:ToGroupID;references:ID"`

	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at" gorm:"index"`
}
