package controllers

import (
	"encoding/json"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"nmteasy_backend/models"
	"nmteasy_backend/utils"
)

func SignUp(w http.ResponseWriter, r *http.Request) {
	var user models.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "failed to decode the body")
		return
	}

	if !utils.NameRegex.MatchString(user.LastName) || !utils.NameRegex.MatchString(user.FirstName) {
		utils.RespondWithError(w, http.StatusBadRequest, "name is not correct")
		return
	}

	if !utils.EmailRegex.MatchString(user.Email) {
		utils.RespondWithError(w, http.StatusBadRequest, "email is not correct")
		return
	}

	if len(user.Password) < 8 {
		utils.RespondWithError(w, http.StatusBadRequest, "password must be 8 characters or more")
		return
	}

	if !utils.IsValidPassword(user.Password) {
		utils.RespondWithError(w, http.StatusBadRequest, "password must contain 1 uppercase, 1 lowercase, 1 number")
		return
	}

	var dbuser models.User
	models.DB.Where("email = ?", user.Email).First(&dbuser)

	//checks if email is already registered or not
	if dbuser.Email != "" {
		utils.RespondWithError(w, http.StatusBadRequest, "user already exists")
		return
	}

	models.DB.Where("username = ?", user.Username).First(&dbuser)

	//checks if username is already registered or not
	if dbuser.Username != "" {
		utils.RespondWithError(w, http.StatusBadRequest, "username already exists")
		return
	}

	user.Password, err = utils.GenerateHashPassword(user.Password)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "error in password hash")
		return
	}

	//insert user details in database
	user.ID = uuid.New()
	models.DB.Create(&user)

	token, err := utils.GenerateJWT(user.Email)
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "failed to generate the token")
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, map[string]interface{}{"token": token, "user": &user})
}

func SignIn(w http.ResponseWriter, r *http.Request) {
	var authDetails models.Authentication
	err := json.NewDecoder(r.Body).Decode(&authDetails)
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "failed to decode the body")
		return
	}

	var authUser models.User
	models.DB.Where("email = ?", authDetails.Email).First(&authUser)
	if authUser.Email == "" {
		utils.RespondWithError(w, http.StatusBadRequest, "no user found")
		return
	}

	check := checkPasswordHash(authDetails.Password, authUser.Password)

	if !check {
		utils.RespondWithError(w, http.StatusBadRequest, "email or password is incorrect")
		return
	}

	token, err := utils.GenerateJWT(authUser.Email)
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "failed to generate the token")
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, token)
}

func checkPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
