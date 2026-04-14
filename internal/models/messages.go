package models

import "time"

type Message struct {
	ID      string `json:"id" gorm:"primaryKey"`
	Content string `json:"content"`

	FromID   string `json:"from_id" gorm:"column:from_user;index;not null"`
	FromUser User   `json:"from_user" gorm:"foreignKey:FromID;references:ID"`

	ToUserID *string `json:"to_user_id,omitempty" gorm:"index"`
	ToUser   *User   `json:"to_user,omitempty" gorm:"foreignKey:ToUserID;references:ID"`

	ToGroupID *string `json:"to_group_id,omitempty" gorm:"index"`
	ToGroup   *Group  `json:"to_group,omitempty" gorm:"foreignKey:ToGroupID;references:ID"`

	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at" gorm:"index"`
}
