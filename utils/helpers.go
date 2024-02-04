package utils

import (
	"encoding/json"
	"fmt"
	"net/http"
	"nmteasy_backend/models"
	"strings"
)

func RespondWithError(w http.ResponseWriter, code int, message string) {
	respondWithJSON(w, code, map[string]string{"error": message})
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
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
