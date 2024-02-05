package utils

import (
	"regexp"
	"unicode"
)

func IsValidPassword(s string) bool {
	var (
		hasUpper  = false
		hasLower  = false
		hasNumber = false
	)

	for _, char := range s {
		switch {
		case unicode.IsUpper(char):
			hasUpper = true
		case unicode.IsLower(char):
			hasLower = true
		case unicode.IsNumber(char):
			hasNumber = true
		}
	}
	return hasUpper && hasLower && hasNumber
}

var EmailRegex = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")
var NameRegex = regexp.MustCompile("^[^_]{1,20}$")
var UsernameRegex = regexp.MustCompile("^[^-][^ ]{1,20}$")
