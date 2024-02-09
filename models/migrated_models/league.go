package migrated_models

import (
	"github.com/google/uuid"
)

type League struct {
	ID     uuid.UUID `json:"id" gorm:"primary_key"`
	Title  string    `json:"title"`
	Points int       `gorm:"not null" json:"points"`

	Users []User `json:"users"`
}
