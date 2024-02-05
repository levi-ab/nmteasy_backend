package controllers

import (
	"github.com/gorilla/mux"
	"net/http"
	"nmteasy_backend/middleware"
)

func New() http.Handler {
	router := mux.NewRouter()

	router.HandleFunc("/history-questions/{lessonID}", middleware.Protected(GetHistoryQuestions)).Methods("GET")
	router.HandleFunc("/history-question-explanation/{questionID}", middleware.Protected(GetHistoryQuestionExplanation)).Methods("GET")
	router.HandleFunc("/history-lessons", middleware.Protected(GetHistoryLessons)).Methods("GET")
	router.HandleFunc("/sign-up", SignUp).Methods("POST")
	router.HandleFunc("/sign-in", SignIn).Methods("POST")
	router.HandleFunc("/user/edit", EditUser).Methods("POST")

	return router
}
