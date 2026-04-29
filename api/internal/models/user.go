package models

import "time"

type UserCreateDTO struct {
	Name     string `json:"name" validate:"required"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
}

type UserLoginDTO struct {
	Name     string `json:"name" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type UserUpdateDTO struct {
	Name  string `json:"name" validate:"required"`
	Email string `json:"email" validate:"required,email"`
}

type User struct {
	ID        string `gorm:"primaryKey;type:uuid" json:"id"`
	Name      string `gorm:"uniqueIndex" json:"name"`
	Email     string `gorm:"uniqueIndex" json:"email"`
	Password  string `json:"-"` // don't return password in json
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time
}

type UserContactDTO struct {
	ContactID string `json:"contact_id" validate:"required"`
}
