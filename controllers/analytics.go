package controllers

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"gorm.io/gorm"
	"net/http"
	"nmteasy_backend/models"
	"nmteasy_backend/models/migrated_models"
	"nmteasy_backend/utils"
	"time"
)

func AddLessonAnalytics(w http.ResponseWriter, r *http.Request) {
	user := utils.GetCurrentUser(r)
	if user == nil {
		utils.RespondWithError(w, http.StatusForbidden, "wrong token")
		return
	}

	lessonType := mux.Vars(r)["lessonType"]

	type FullLessonAnalytic struct {
		LessonAnalytic struct {
			LessonID          uuid.UUID `json:"lesson_id"`
			TimeSpent         int       `json:"time_spent"`
			RightAnswersCount int       `json:"right_answers_count"`
			QuestionsCount    int       `json:"questions_count"`
		} `json:"lesson_analytic"`
		QuestionAnalytics []struct {
			QuestionID    uuid.UUID `json:"question_id"`
			TimeSpent     int       `json:"time_spent"`
			AnsweredRight bool      `json:"answered_right"`
		} `json:"question_analytics"`
	}

	var model FullLessonAnalytic

	err := json.NewDecoder(r.Body).Decode(&model)
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "failed to decode the body")
		return
	}

	var existingLessonAnalytic models.LessonAnalytic

	query := fmt.Sprintf(`SELECT id, user_id, right_answers_count,questions_count, time_spent,  created_at, updated_at, %s_lesson_id as lesson_id
								FROM %s_lesson_analytics WHERE %s_lesson_id = ? AND user_id = ?`, lessonType, lessonType, lessonType)

	if err := models.DB.Raw(query, model.LessonAnalytic.LessonID, user.ID).Scan(&existingLessonAnalytic).Error; err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			utils.RespondWithError(w, http.StatusInternalServerError, "Failed to get explanation")
			return
		}
	}

	now := time.Now().Format("2006-01-02 15:04:05.999 -0700")

	user.Points = user.Points + model.LessonAnalytic.RightAnswersCount
	user.WeeklyPoints = user.WeeklyPoints + model.LessonAnalytic.RightAnswersCount

	if err = models.DB.Save(&user).Error; err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "failed updated the user points")
		return
	}

	if existingLessonAnalytic.ID == uuid.Nil {
		IDToInsert := uuid.New()
		query := fmt.Sprintf(`INSERT INTO %s_lesson_analytics (id, %s_lesson_id, user_id, right_answers_count, questions_count, time_spent, created_at, updated_at)
					VALUES (?, ?, ?, ?, ?, ?, ?, ?);`, lessonType, lessonType)
		if err := models.DB.Exec(query, IDToInsert, model.LessonAnalytic.LessonID, user.ID, model.LessonAnalytic.RightAnswersCount, model.LessonAnalytic.QuestionsCount, model.LessonAnalytic.TimeSpent, now, now).Error; err != nil {
			utils.RespondWithError(w, http.StatusInternalServerError, "failed to save the analytic for lesson")
			return
		}
	} else {
		if model.LessonAnalytic.RightAnswersCount > existingLessonAnalytic.RightAnswersCount {
			existingLessonAnalytic.RightAnswersCount = model.LessonAnalytic.RightAnswersCount
		}

		query := fmt.Sprintf(`UPDATE %s_lesson_analytics
				SET
					right_answers_count = ?,
					questions_count = ?,
					time_spent = ?,
				    updated_at = ?
				WHERE
    id = ?;`, lessonType)

		if err := models.DB.Exec(query, existingLessonAnalytic.RightAnswersCount, model.LessonAnalytic.QuestionsCount, model.LessonAnalytic.TimeSpent, now, existingLessonAnalytic.ID).Error; err != nil {
			utils.RespondWithError(w, http.StatusInternalServerError, "failed to save the analytic for lesson")
			return
		}
	}

	var questionIDS []uuid.UUID
	for _, qa := range model.QuestionAnalytics {
		questionIDS = append(questionIDS, qa.QuestionID)
	}

	if err = models.DB.Exec(fmt.Sprintf("DELETE FROM %s_question_analytics WHERE user_id = ? AND %s_question_id IN ? ", lessonType, lessonType), user.ID, questionIDS).Error; err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "failed to remove old analytics for questions")
		return
	}

	if lessonType == "history" { //todo need to make it generic somehow
		var questionsAnalytics []migrated_models.HistoryQuestionAnalytic

		for _, questionAnalytic := range model.QuestionAnalytics {
			questionsAnalytics = append(questionsAnalytics, migrated_models.HistoryQuestionAnalytic{
				ID:                uuid.New(),
				HistoryQuestionID: questionAnalytic.QuestionID,
				UserID:            user.ID,
				AnsweredRight:     questionAnalytic.AnsweredRight,
				TimeSpent:         questionAnalytic.TimeSpent,
			})
		}

		if err = models.DB.Save(&questionsAnalytics).Error; err != nil {
			utils.RespondWithError(w, http.StatusInternalServerError, "failed to save analytics for questions")
			return
		}

		utils.RespondWithJSON(w, http.StatusOK, nil)
		return
	}

	if lessonType == "ukrainian" {
		var questionsAnalytics []migrated_models.UkrainianQuestionAnalytic

		for _, questionAnalytic := range model.QuestionAnalytics {
			questionsAnalytics = append(questionsAnalytics, migrated_models.UkrainianQuestionAnalytic{
				ID:                  uuid.New(),
				UkrainianQuestionID: questionAnalytic.QuestionID,
				UserID:              user.ID,
				AnsweredRight:       questionAnalytic.AnsweredRight,
				TimeSpent:           questionAnalytic.TimeSpent,
			})
		}

		if err = models.DB.Save(&questionsAnalytics).Error; err != nil {
			utils.RespondWithError(w, http.StatusInternalServerError, "failed to save analytics for questions")
			return
		}

		utils.RespondWithJSON(w, http.StatusOK, nil)
		return
	}

	var questionsAnalytics []migrated_models.BiologyQuestionAnalytic

	for _, questionAnalytic := range model.QuestionAnalytics {
		questionsAnalytics = append(questionsAnalytics, migrated_models.BiologyQuestionAnalytic{
			ID:                uuid.New(),
			BiologyQuestionID: questionAnalytic.QuestionID,
			UserID:            user.ID,
			AnsweredRight:     questionAnalytic.AnsweredRight,
			TimeSpent:         questionAnalytic.TimeSpent,
		})
	}

	if err = models.DB.Save(&questionsAnalytics).Error; err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "failed to save analytics for questions")
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, nil)
	return
}

func GetWeeklyQuestionAnalytics(w http.ResponseWriter, r *http.Request) {
	user := utils.GetCurrentUser(r)
	if user == nil {
		utils.RespondWithError(w, http.StatusForbidden, "wrong token")
		return
	}

	lessonType := mux.Vars(r)["lessonType"]

	var analytics []struct {
		DayOfWeek string `gorm:"day_of_week"`
		Count     int    `gorm:"count"`
	}

	if err := models.DB.Raw(fmt.Sprintf(`
        SELECT 
    date_trunc('day', created_at) AS day_of_week, 
     SUM(COUNT(*)) OVER (ORDER BY date_trunc('day', created_at)) AS count
FROM %s_question_analytics
WHERE 
    user_id = ? AND
    created_at >= current_date - interval '7 days'
GROUP BY day_of_week
ORDER BY day_of_week;
    `, lessonType), user.ID).Scan(&analytics).Error; err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "failed to get analytics")
		return
	}

	// Convert the result to a map for easy access
	result := make(map[string]int)
	for _, entry := range analytics {
		result[entry.DayOfWeek] = entry.Count
	}

	utils.RespondWithJSON(w, http.StatusOK, result)
}
