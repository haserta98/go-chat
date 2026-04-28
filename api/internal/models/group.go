package models

import "time"

type Group struct {
	ID   string `json:"id" gorm:"primaryKey;type:uuid"`
	Name string `json:"name"`

	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at" gorm:"index"`
}

type GroupMember struct {
	ID string `json:"id" gorm:"primaryKey;type:uuid"`

	GroupID string `json:"group_id" gorm:"index;not null"`
	Group   Group  `json:"group" gorm:"foreignKey:GroupID;references:ID"`
	UserID  string `json:"user_id" gorm:"index;not null"`
	User    User   `json:"user" gorm:"foreignKey:UserID;references:ID"`
	Role    string `json:"role" gorm:"type:varchar(20);default:'member'"`

	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at" gorm:"index"`
}
