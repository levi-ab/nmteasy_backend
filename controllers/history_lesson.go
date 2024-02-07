package controllers

import (
	"github.com/gorilla/mux"
	"net/http"
	"nmteasy_backend/models"
	"nmteasy_backend/utils"
	"sort"
)

func GetHistoryQuestions(w http.ResponseWriter, r *http.Request) {
	var historyQuestions []models.HistoryQuestion
	paramValue := mux.Vars(r)["lessonID"]

	if err := models.DB.Where("history_lesson_id = ? AND topic != 'error'", paramValue).Find(&historyQuestions).Error; err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Failed to get question")
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, historyQuestions)
}

func GetHistoryLessons(w http.ResponseWriter, r *http.Request) {
	user := utils.GetCurrentUser(r)
	if user == nil {
		utils.RespondWithError(w, http.StatusForbidden, "wrong token")
		return
	}

	var historyLessons []models.HistoryLessonWithProperTitle

	query := `
    SELECT h.id, h.title, h.created_at, a.id as analytic_id, a.user_id, COALESCE(a.right_answers_count, 0), COALESCE(a.questions_count, 0),  COALESCE(a.created_at, '0001-01-01'::timestamp) as analytic_created_at,  COALESCE(a.updated_at, '0001-01-01'::timestamp) as analytic_updated_at, COALESCE(a.time_spent,0) as time_spent
    FROM history_lessons h
    LEFT JOIN history_lesson_analytics a ON h.id = a.history_lesson_id AND a.user_id = ?
    ORDER BY h.created_at, h.title
`

	rows, err := models.DB.Raw(query, user.ID).Rows()
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Failed to execute query")
		return
	}
	defer rows.Close()

	for rows.Next() {
		var historyLesson models.HistoryLessonWithProperTitle
		var analytic models.HistoryLessonAnalytic

		if err := rows.Scan(
			&historyLesson.ID,
			&historyLesson.Title,
			&historyLesson.CreatedAt,

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

		historyLesson.HistoryLessonAnalytic = analytic
		historyLessons = append(historyLessons, historyLesson)
	}

	if err := rows.Err(); err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Error iterating over rows")
		return
	}
	groupedLessonsMap := make(map[string][]models.HistoryLessonWithProperTitle)

	for _, historyLesson := range historyLessons {
		properTitle, generalTitle := utils.FormatLessonTopic(historyLesson.Title)

		if lessons, ok := groupedLessonsMap[generalTitle]; ok {
			groupedLessonsMap[generalTitle] = append(lessons, models.HistoryLessonWithProperTitle{
				Title:                 historyLesson.Title,
				ID:                    historyLesson.ID,
				ProperTitle:           properTitle,
				CreatedAt:             historyLesson.CreatedAt,
				HistoryLessonAnalytic: historyLesson.HistoryLessonAnalytic,
			})
		} else {
			groupedLessonsMap[generalTitle] = []models.HistoryLessonWithProperTitle{{
				Title:                 historyLesson.Title,
				ID:                    historyLesson.ID,
				ProperTitle:           properTitle,
				CreatedAt:             historyLesson.CreatedAt,
				HistoryLessonAnalytic: historyLesson.HistoryLessonAnalytic,
			}}
		}
	}

	type GroupedLesson struct {
		Title string
		Data  []models.HistoryLessonWithProperTitle
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

func GetHistoryQuestionExplanation(w http.ResponseWriter, r *http.Request) {
	var explanation models.HistoryQuestionExplanation
	paramValue := mux.Vars(r)["questionID"]

	if err := models.DB.Where("history_question_id = ?", paramValue).Find(&explanation).Error; err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Failed to get question")
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, explanation)
}
