package migrated_models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	ID           uuid.UUID `json:"id" gorm:"primary_key"`
	FirstName    string    `json:"first_name"`
	LastName     string    `json:"last_name"`
	Username     string    `json:"username"`
	Email        string    `gorm:"unique" json:"email"`
	Points       int       `gorm:"not null" json:"points"`
	WeeklyPoints int       `gorm:"not null" json:"weekly-points"`
	Password     string    `json:"password"`

	LeagueID *uuid.UUID `json:"league_id"`
	League   League     `json:"league"`
}
