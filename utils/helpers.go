package utils

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"nmteasy_backend/common"
	"nmteasy_backend/models"
	"strings"
	"time"
)

func RespondWithError(w http.ResponseWriter, code int, message string) {
	RespondWithJSON(w, code, map[string]string{"error": message})
}

func RespondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}

func FormatLessonTopic(input string) (string, string) {
	// Split the input string into topic and number
	parts := strings.Split(input, "#")
	if len(parts) != 2 {
		return "", ""
	}
	topicKey := parts[0]
	number := parts[1]

	topic, ok := models.HistoryTopics[topicKey]
	if !ok {
		return input, "" // Return the original input if the key is not found
	}

	// Format the result
	result := fmt.Sprintf("%s - %s", topic, number)
	return result, topic
}

func GenerateHashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func GenerateJWT(email string) (string, error) {
	var mySigningKey = []byte(common.SECRET_KEY)
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)

	claims["authorized"] = true
	claims["email"] = email
	claims["exp"] = time.Now().Add(365 * 24 * time.Hour).Unix()

	tokenString, err := token.SignedString(mySigningKey)

	if err != nil {
		fmt.Errorf("something Went Wrong: %s", err.Error())
		return "", err
	}
	return tokenString, nil
}

func GetCurrentUser(r *http.Request) *models.User {
	token, err := GetToken(r)
	if err != nil || token == "" {
		return nil
	}

	mySigningKey := []byte(common.SECRET_KEY)
	parsedToken, err := parseToken(token, mySigningKey)
	if err != nil || parsedToken == nil {
		return nil
	}

	if claims, ok := parsedToken.Claims.(jwt.MapClaims); ok {
		if email, ok := claims["email"].(string); ok && email != "" {
			user, err := getUserByEmail(email)
			if err == nil {
				return user
			}
		}
	}

	return nil
}

func parseToken(tokenString string, signingKey []byte) (*jwt.Token, error) {
	return jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if method, ok := token.Method.(*jwt.SigningMethodHMAC); !ok || method != jwt.SigningMethodHS256 {
			return nil, fmt.Errorf("invalid token signing method")
		}
		return signingKey, nil
	})
}

func getUserByEmail(email string) (*models.User, error) {
	var user models.User
	err := models.DB.Where("email = ?", email).Find(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func GetToken(r *http.Request) (string, error) {
	h := r.Header.Get("Authorization")
	if h == "" || !strings.Contains(h, "Bearer ") {
		return "", errors.New("missing token")
	}
	return strings.TrimPrefix(h, "Bearer "), nil
}
