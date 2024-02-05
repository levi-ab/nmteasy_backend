package controllers

import (
	"encoding/json"
	"net/http"
	"nmteasy_backend/models"
	"nmteasy_backend/utils"
)

func EditUser(w http.ResponseWriter, r *http.Request) {
	var userModelPayload models.User
	currentUser := utils.GetCurrentUser(r)
	if currentUser == nil {
		utils.RespondWithError(w, http.StatusForbidden, "no token found")
		return
	}

	err := json.NewDecoder(r.Body).Decode(&userModelPayload)
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "failed to decode the body")
		return
	}

	if !utils.NameRegex.MatchString(userModelPayload.LastName) || !utils.NameRegex.MatchString(userModelPayload.FirstName) {
		utils.RespondWithError(w, http.StatusBadRequest, "name is not correct")
		return
	}

	if !utils.EmailRegex.MatchString(userModelPayload.Email) {
		utils.RespondWithError(w, http.StatusBadRequest, "email is not correct")
		return
	}

	if !utils.UsernameRegex.MatchString(userModelPayload.Username) {
		utils.RespondWithError(w, http.StatusBadRequest, "username is not correct")
		return
	}

	if userModelPayload.Email != currentUser.Email {
		var dbuser models.User
		models.DB.Where("email = ?", userModelPayload.Email).First(&dbuser)

		//checks if email is already registered or not
		if dbuser.Email != "" {
			utils.RespondWithError(w, http.StatusBadRequest, "email already exists")
			return
		}

		currentUser.Email = userModelPayload.Email
	}

	if userModelPayload.Username != currentUser.Username {
		var dbuser models.User
		models.DB.Where("username = ?", userModelPayload.Username).First(&dbuser)

		//checks if username is already registered or not
		if dbuser.Username != "" {
			utils.RespondWithError(w, http.StatusBadRequest, "username already exists")
			return
		}

		currentUser.Username = userModelPayload.Username
	}

	currentUser.FirstName = userModelPayload.FirstName
	currentUser.LastName = userModelPayload.LastName

	if err := models.DB.Save(&currentUser).Error; err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "failed to save the user")
		return
	}

	token, err := utils.GenerateJWT(userModelPayload.Email)
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "failed to generate the token")
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, map[string]interface{}{"token": token, "user": &currentUser})
}
