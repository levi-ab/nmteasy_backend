package controllers

import (
	"github.com/gorilla/mux"
	"net/http"
	"nmteasy_backend/middleware"
	"nmteasy_backend/sockets"
)

func New() http.Handler {
	router := mux.NewRouter()

	//lessons
	router.HandleFunc("/questions/{lessonType}/{lessonID}", middleware.Protected(GetQuestionsByLesson)).Methods("GET")
	router.HandleFunc("/question-explanation/{lessonType}/{questionID}", middleware.Protected(GetQuestionExplanation)).Methods("GET")
	router.HandleFunc("/lessons/{lessonType}", middleware.Protected(GetLessons)).Methods("GET")

	//analytics
	router.HandleFunc("/lessons-analytic/{lessonType}", middleware.Protected(AddLessonAnalytics)).Methods("POST")
	router.HandleFunc("/questions-analytic/{lessonType}", middleware.Protected(AddQuestionsAnalytics)).Methods("POST")

	router.HandleFunc("/weekly-analytics/{lessonType}", middleware.Protected(GetWeeklyQuestionAnalytics)).Methods("GET")

	router.HandleFunc("/lesson-complaint", middleware.Protected(AddComplaint)).Methods("POST")

	//user
	router.HandleFunc("/sign-up", SignUp).Methods("POST")
	router.HandleFunc("/sign-in", SignIn).Methods("POST")
	router.HandleFunc("/user/edit", middleware.Protected(EditUser)).Methods("POST")

	//leagues
	router.HandleFunc("/leagues", middleware.Protected(GetLeagues)).Methods("GET")
	router.HandleFunc("/current-league", middleware.Protected(GetCurrentLeague)).Methods("GET")

	//questions
	router.HandleFunc("/random-questions/{lessonType}", middleware.Protected(GetRandomQuestions)).Methods("GET")
	router.HandleFunc("/match-questions/{lessonType}", middleware.Protected(GetMatchQuestions)).Methods("GET")
	router.HandleFunc("/wrong-answer-questions/{lessonType}", middleware.Protected(GetWrongAnsweredQuestions)).Methods("GET")
	router.HandleFunc("/image-questions/{lessonType}", middleware.Protected(GetImageQuestions)).Methods("GET")

	router.HandleFunc("/shop/products", middleware.Protected(GetProducts)).Methods("GET")
	router.HandleFunc("/shop/purchase", middleware.Protected(PurchaseProduct)).Methods("POST")
	router.HandleFunc("/shop/inventory", middleware.Protected(GetUserInventory)).Methods("GET")
	router.HandleFunc("/shop/activate-skin", middleware.Protected(ActivateSkin)).Methods("POST")

	hub := sockets.NewHub()
	go hub.Run()
	router.HandleFunc("/ws/{token}", func(w http.ResponseWriter, r *http.Request) {
		sockets.ServeWs(hub, w, r)
	})

	return router
}
