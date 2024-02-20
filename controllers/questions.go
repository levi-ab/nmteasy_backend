package controllers

import (
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
	"nmteasy_backend/models"
	"nmteasy_backend/utils"
)

func GetRandomQuestions(w http.ResponseWriter, r *http.Request) {
	user := utils.GetCurrentUser(r)
	if user == nil {
		utils.RespondWithError(w, http.StatusForbidden, "wrong token")
		return
	}

	lessonType := mux.Vars(r)["lessonType"]
	var questions []models.Question

	rows, err := models.DB.
		Raw(fmt.Sprintf(`SELECT id, question_text, question_image, answers, type, right_answer, topic, created_at, updated_at, %s_lesson_id as lesson_id
    FROM %s_questions  ORDER BY RANDOM()  LIMIT 20`, lessonType, lessonType)).Rows()

	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Failed to get questions")
		return
	}
	defer rows.Close()

	for rows.Next() {
		var question models.Question
		if err := models.DB.ScanRows(rows, &question); err != nil {
			utils.RespondWithError(w, http.StatusInternalServerError, "Failed to scan rows")
			return
		}

		questions = append(questions, question)
	}

	if err := rows.Err(); err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Error during rows iteration")
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, questions)

}

func GetMatchQuestions(w http.ResponseWriter, r *http.Request) {
	user := utils.GetCurrentUser(r)
	if user == nil {
		utils.RespondWithError(w, http.StatusForbidden, "wrong token")
		return
	}

	lessonType := mux.Vars(r)["lessonType"]
	var questions []models.Question

	rows, err := models.DB.
		Raw(fmt.Sprintf(`SELECT q.id, q.question_text, q.question_image, q.answers, q.type, q.right_answer, q.topic, q.created_at, q.updated_at, q.%s_lesson_id as lesson_id
    FROM %s_questions as q
    LEFT JOIN %s_question_analytics as an ON q.id = an.%s_question_id AND an.answered_right = 'true'
    WHERE an.%s_question_id IS NULL AND q.type = 'match_two_answers_rows'
    ORDER BY q.created_at
    LIMIT 20`, lessonType, lessonType, lessonType, lessonType, lessonType)).Rows()

	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Failed to get questions")
		return
	}
	defer rows.Close()

	for rows.Next() {
		var question models.Question
		if err := models.DB.ScanRows(rows, &question); err != nil {
			utils.RespondWithError(w, http.StatusInternalServerError, "Failed to scan rows")
			return
		}

		questions = append(questions, question)
	}

	if err := rows.Err(); err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Error during rows iteration")
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, questions)

}

func GetWrongAnsweredQuestions(w http.ResponseWriter, r *http.Request) {
	user := utils.GetCurrentUser(r)
	if user == nil {
		utils.RespondWithError(w, http.StatusForbidden, "wrong token")
		return
	}

	lessonType := mux.Vars(r)["lessonType"]
	var questions []models.Question

	rows, err := models.DB.
		Raw(fmt.Sprintf(`SELECT q.id, q.question_text, q.question_image, q.answers, q.type, q.right_answer, q.topic, q.created_at, q.updated_at, q.%s_lesson_id as lesson_id
    FROM %s_questions as q
    JOIN %s_question_analytics as an ON q.id = an.%s_question_id AND an.answered_right = 'false'
    ORDER BY an.created_at DESC 	
    LIMIT 20`, lessonType, lessonType, lessonType, lessonType)).Rows()

	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Failed to get questions")
		return
	}
	defer rows.Close()

	for rows.Next() {
		var question models.Question
		if err := models.DB.ScanRows(rows, &question); err != nil {
			utils.RespondWithError(w, http.StatusInternalServerError, "Failed to scan rows")
			return
		}

		questions = append(questions, question)
	}

	if err := rows.Err(); err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Error during rows iteration")
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, questions)

}
