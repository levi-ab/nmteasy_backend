package sockets

import (
	"github.com/google/uuid"
	"nmteasy_backend/models"
	"nmteasy_backend/models/migrated_models"
)

func UpdateUserPoints(userID uuid.UUID, points int) {
	var user migrated_models.User
	err := models.DB.Where("id = ?", userID).First(&user).Error
	if err == nil {
		user.Points += points
		models.DB.Save(&user)
	}
}
