package controllers

import (
	"encoding/json"
	"github.com/google/uuid"
	"net/http"
	"nmteasy_backend/models"
	"nmteasy_backend/utils"
)

func AddHistoryLessonAnalytics(w http.ResponseWriter, r *http.Request) {
	user := utils.GetCurrentUser(r)
	var analyticToInsert models.HistoryLessonAnalytic
	if user == nil {
		utils.RespondWithError(w, http.StatusForbidden, "wrong token")
		return
	}

	err := json.NewDecoder(r.Body).Decode(&analyticToInsert)
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "failed to decode the body")
		return
	}

	var existingHistoryLessonAnalytic models.HistoryLessonAnalytic
	err = models.DB.Where("user_id = ? AND history_lesson_id = ?", user.ID, analyticToInsert.HistoryLessonID).Find(&existingHistoryLessonAnalytic).Error
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "failed to get the analytic")
		return
	}

	if existingHistoryLessonAnalytic.ID == uuid.Nil {
		analyticToInsert.ID = uuid.New()
		analyticToInsert.UserID = user.ID
		err = models.DB.Save(&analyticToInsert).Error
		if err != nil {
			utils.RespondWithError(w, http.StatusInternalServerError, "failed to save the analytic")
			return
		}

		utils.RespondWithJSON(w, http.StatusOK, nil)
		return
	}

	if analyticToInsert.RightAnswersCount > existingHistoryLessonAnalytic.RightAnswersCount {
		existingHistoryLessonAnalytic.RightAnswersCount = analyticToInsert.RightAnswersCount
	}

	existingHistoryLessonAnalytic.TimeSpent = analyticToInsert.TimeSpent

	err = models.DB.Save(&existingHistoryLessonAnalytic).Error
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "failed to save the analytic")
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, nil)
}
