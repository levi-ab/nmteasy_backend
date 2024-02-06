package controllers

import (
	"encoding/json"
	"github.com/google/uuid"
	"net/http"
	"nmteasy_backend/models"
	"nmteasy_backend/utils"
)

func AddComplaint(w http.ResponseWriter, r *http.Request) {
	user := utils.GetCurrentUser(r)
	var complaintToInsert models.Complaint
	if user == nil {
		utils.RespondWithError(w, http.StatusForbidden, "wrong token")
		return
	}

	err := json.NewDecoder(r.Body).Decode(&complaintToInsert)
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "failed to decode the body")
		return
	}

	complaintToInsert.UserID = user.ID
	complaintToInsert.ID = uuid.New()

	err = models.DB.Save(&complaintToInsert).Error
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "failed to save the complaint")
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, nil)
}
