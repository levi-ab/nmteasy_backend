package controllers

import (
	"github.com/gorilla/mux"
	"net/http"
)

func New() http.Handler {
	router := mux.NewRouter()

	router.HandleFunc("/history-questions/{lessonID}", GetHistoryQuestions).Methods("GET")
	router.HandleFunc("/history-lessons", GetHistoryLessons).Methods("GET")

	return router
}
