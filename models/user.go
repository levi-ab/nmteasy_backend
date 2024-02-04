package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	ID        uuid.UUID `json:"id" gorm:"primary_key"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	Username  string    `json:"username"`
	Email     string    `gorm:"unique" json:"email"`
	Password  string    `json:"password"`
}

type Authentication struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
