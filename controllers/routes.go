package controllers

import (
	"github.com/gorilla/mux"
	"net/http"
	"nmteasy_backend/middleware"
)

func New() http.Handler {
	router := mux.NewRouter()

	router.HandleFunc("/questions/{lessonType}/{lessonID}", middleware.Protected(GetQuestionsByLesson)).Methods("GET")
	router.HandleFunc("/question-explanation/{lessonType}/{questionID}", middleware.Protected(GetQuestionExplanation)).Methods("GET")
	router.HandleFunc("/lessons/{lessonType}", middleware.Protected(GetLessons)).Methods("GET")

	//analytics
	router.HandleFunc("/lessons-analytic/{lessonType}", middleware.Protected(AddLessonAnalytics)).Methods("POST")
	router.HandleFunc("/weekly-analytics/{lessonType}", middleware.Protected(GetWeeklyQuestionAnalytics)).Methods("GET")

	router.HandleFunc("/lesson-complaint", middleware.Protected(AddComplaint)).Methods("POST")

	//user
	router.HandleFunc("/sign-up", SignUp).Methods("POST")
	router.HandleFunc("/sign-in", SignIn).Methods("POST")
	router.HandleFunc("/user/edit", middleware.Protected(EditUser)).Methods("POST")

	return router
}
