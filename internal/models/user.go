package models

import "time"

type UserCreateDTO struct {
	Name  string `json:"name" validate:"required"`
	Email string `json:"email" validate:"required,email"`
}

type UserUpdateDTO struct {
	Name  string `json:"name" validate:"required"`
	Email string `json:"email" validate:"required,email"`
}

type User struct {
	ID        uint `gorm:"primaryKey" json:"id"`
	Name      string
	Email     string
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time
}
