package controllers

import (
	"errors"
	"fmt"
	"github.com/gorilla/mux"
	"gorm.io/gorm"
	"net/http"
	"nmteasy_backend/models"
	"nmteasy_backend/utils"
	"sort"
)

func GetQuestionsByLesson(w http.ResponseWriter, r *http.Request) {
	user := utils.GetCurrentUser(r)
	if user == nil {
		utils.RespondWithError(w, http.StatusForbidden, "wrong token")
		return
	}

	lessonType := mux.Vars(r)["lessonType"]
	lessonID := mux.Vars(r)["lessonID"]
	var questions []models.Question

	rows, err := models.DB.
		Raw(fmt.Sprintf(`SELECT id, question_text, question_image, answers, type, right_answer, topic, created_at, updated_at, %s_lesson_id as lesson_id
								FROM %s_questions WHERE %s_lesson_id = ? AND topic != 'error'`, lessonType, lessonType, lessonType), lessonID).Rows()

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

func GetLessons(w http.ResponseWriter, r *http.Request) {
	user := utils.GetCurrentUser(r)
	if user == nil {
		utils.RespondWithError(w, http.StatusForbidden, "wrong token")
		return
	}

	lessonType := mux.Vars(r)["lessonType"]

	var lessons []models.LessonWithProperTitle

	query := fmt.Sprintf(`
    SELECT h.id, h.title, h.created_at, a.id as analytic_id, a.user_id, COALESCE(a.right_answers_count, 0), COALESCE(a.questions_count, 0), COALESCE(a.created_at, '0001-01-01'::timestamp) as analytic_created_at, COALESCE(a.updated_at, '0001-01-01'::timestamp) as analytic_updated_at, COALESCE(a.time_spent,0) as time_spent
    FROM %s_lessons h
    LEFT JOIN %s_lesson_analytics a ON h.id = a.%s_lesson_id AND a.user_id = ?
    ORDER BY h.created_at, h.title
`, lessonType, lessonType, lessonType)

	rows, err := models.DB.Raw(query, user.ID).Rows()
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Failed to execute query")
		return
	}
	defer rows.Close()

	for rows.Next() {
		var lesson models.LessonWithProperTitle
		var analytic models.LessonAnalytic

		if err := rows.Scan(
			&lesson.ID,
			&lesson.Title,
			&lesson.CreatedAt,

			&analytic.ID,
			&analytic.UserID,
			&analytic.RightAnswersCount,
			&analytic.QuestionsCount,
			&analytic.CreatedAt,
			&analytic.UpdatedAt,
			&analytic.TimeSpent,
		); err != nil {
			utils.RespondWithError(w, http.StatusInternalServerError, "Failed to scan rows")
			return
		}

		lesson.LessonAnalytic = analytic
		lessons = append(lessons, lesson)
	}

	if err := rows.Err(); err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Error iterating over rows")
		return
	}
	groupedLessonsMap := make(map[string][]models.LessonWithProperTitle)

	for _, lesson := range lessons {
		properTitle, generalTitle := utils.FormatLessonTopic(lesson.Title)

		if lessons, ok := groupedLessonsMap[generalTitle]; ok {
			groupedLessonsMap[generalTitle] = append(lessons, models.LessonWithProperTitle{
				Title:          lesson.Title,
				ID:             lesson.ID,
				ProperTitle:    properTitle,
				CreatedAt:      lesson.CreatedAt,
				LessonAnalytic: lesson.LessonAnalytic,
			})
		} else {
			groupedLessonsMap[generalTitle] = []models.LessonWithProperTitle{{
				Title:          lesson.Title,
				ID:             lesson.ID,
				ProperTitle:    properTitle,
				CreatedAt:      lesson.CreatedAt,
				LessonAnalytic: lesson.LessonAnalytic,
			}}
		}
	}

	type GroupedLesson struct {
		Title string
		Data  []models.LessonWithProperTitle
	}
	//need this cause maps mess up the order of the lessons

	var groupedLessons []GroupedLesson

	for title, lessons := range groupedLessonsMap {
		groupedLessons = append(groupedLessons, GroupedLesson{
			Title: title,
			Data:  lessons,
		})
	}

	sort.Slice(groupedLessons, func(i, j int) bool {
		return groupedLessons[i].Data[0].CreatedAt.Before(groupedLessons[j].Data[0].CreatedAt)
	})

	var result []map[string]interface{}
	for _, group := range groupedLessons {
		result = append(result, map[string]interface{}{
			"title": group.Title,
			"data":  group.Data,
		})
	}

	utils.RespondWithJSON(w, http.StatusOK, result)
}

func GetQuestionExplanation(w http.ResponseWriter, r *http.Request) {
	user := utils.GetCurrentUser(r)
	if user == nil {
		utils.RespondWithError(w, http.StatusForbidden, "wrong token")
		return
	}

	var explanation models.Explanation
	lessonType := mux.Vars(r)["lessonType"]
	questionID := mux.Vars(r)["questionID"]

	query := fmt.Sprintf(`SELECT id, explanation, created_at, updated_at, %s_question_id as question_id
								FROM %s_question_explanations WHERE %s_question_id = ?`, lessonType, lessonType, lessonType)

	if err := models.DB.Raw(query, questionID).Scan(&explanation).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			utils.RespondWithError(w, http.StatusNotFound, "Explanation not found")
		} else {
			utils.RespondWithError(w, http.StatusInternalServerError, "Failed to get explanation")
		}
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, explanation)
}
