package middleware

import (
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt"
	"net/http"
	"nmteasy_backend/common"
	"nmteasy_backend/utils"
	"strings"
)

func Protected(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		token, err := GetToken(r)
		if err != nil || token == "" {
			utils.RespondWithError(w, http.StatusForbidden, "no token found")
			return
		}

		var mySigningKey = []byte(common.SECRET_KEY)

		_, err = jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("there was an error in parsing token")
			}
			return mySigningKey, nil
		})

		if err != nil {
			utils.RespondWithError(w, http.StatusForbidden, "your token has been expired")
			return
		}

		next.ServeHTTP(w, r)
	}
}

func GetToken(r *http.Request) (string, error) {
	h := r.Header.Get("Authorization")
	if h == "" || !strings.Contains(h, "Bearer ") {
		return "", errors.New("missing token")
	}
	return strings.TrimPrefix(h, "Bearer "), nil
}
