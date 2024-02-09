package controllers

import (
	"net/http"
	"nmteasy_backend/models"
	"nmteasy_backend/models/migrated_models"
	"nmteasy_backend/utils"
)

func GetLeagues(w http.ResponseWriter, r *http.Request) {
	user := utils.GetCurrentUser(r)
	if user == nil {
		utils.RespondWithError(w, http.StatusForbidden, "wrong token")
		return
	}

	var leagues []migrated_models.League

	if err := models.DB.Find(&leagues).Error; err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "failed to get leagues")
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, leagues)
}

func GetCurrentLeague(w http.ResponseWriter, r *http.Request) {
	user := utils.GetCurrentUser(r)
	if user == nil {
		utils.RespondWithError(w, http.StatusForbidden, "wrong token")
		return
	}

	var league migrated_models.League

	if err := models.DB.Where("id = ?", user.LeagueID).Preload("Users").First(&league).Error; err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "failed to get leagues")
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, league)
}
