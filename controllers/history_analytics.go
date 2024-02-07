package controllers

import (
	"encoding/json"
	"github.com/google/uuid"
	"net/http"
	"nmteasy_backend/models"
	"nmteasy_backend/utils"
)

func AddHistoryLessonAnalytics(w http.ResponseWriter, r *http.Request) {
	user := utils.GetCurrentUser(r)
	if user == nil {
		utils.RespondWithError(w, http.StatusForbidden, "wrong token")
		return
	}

	type FullHistoryLessonAnalytic struct {
		LessonAnalytic struct {
			HistoryLessonID   uuid.UUID `json:"history_lesson_id"`
			TimeSpent         int       `json:"time_spent"`
			RightAnswersCount int       `json:"right_answers_count"`
			QuestionsCount    int       `json:"questions_count"`
		} `json:"history_lesson_analytic"`
		QuestionAnalytics []struct {
			HistoryQuestionID uuid.UUID `json:"history_question_id"`
			TimeSpent         int       `json:"time_spent"`
			AnsweredRight     bool      `json:"answered_right"`
		} `json:"history_question_analytics"`
	}

	var model FullHistoryLessonAnalytic

	err := json.NewDecoder(r.Body).Decode(&model)
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "failed to decode the body")
		return
	}

	var historyLessonAnalytic models.HistoryLessonAnalytic
	historyLessonAnalytic.HistoryLessonID = model.LessonAnalytic.HistoryLessonID
	historyLessonAnalytic.RightAnswersCount = model.LessonAnalytic.RightAnswersCount
	historyLessonAnalytic.TimeSpent = model.LessonAnalytic.TimeSpent
	historyLessonAnalytic.QuestionsCount = model.LessonAnalytic.QuestionsCount

	var existingHistoryLessonAnalytic models.HistoryLessonAnalytic
	err = models.DB.Where("user_id = ? AND history_lesson_id = ?", user.ID, historyLessonAnalytic.HistoryLessonID).Find(&existingHistoryLessonAnalytic).Error
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "failed to get the analytic for lesson")
		return
	}

	if existingHistoryLessonAnalytic.ID == uuid.Nil {
		historyLessonAnalytic.ID = uuid.New()
		historyLessonAnalytic.UserID = user.ID
		err = models.DB.Save(&historyLessonAnalytic).Error
		if err != nil {
			utils.RespondWithError(w, http.StatusInternalServerError, "failed to save the analytic for lesson")
			return
		}

		utils.RespondWithJSON(w, http.StatusOK, nil)
		return
	}

	if historyLessonAnalytic.RightAnswersCount > existingHistoryLessonAnalytic.RightAnswersCount {
		existingHistoryLessonAnalytic.RightAnswersCount = historyLessonAnalytic.RightAnswersCount
	}

	existingHistoryLessonAnalytic.TimeSpent = historyLessonAnalytic.TimeSpent

	err = models.DB.Save(&existingHistoryLessonAnalytic).Error
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "failed to save the analytic for lesson")
		return
	}

	var questionIDS []uuid.UUID
	for _, qa := range model.QuestionAnalytics {
		questionIDS = append(questionIDS, qa.HistoryQuestionID)
	}

	if err = models.DB.Exec("DELETE FROM history_question_analytics WHERE user_id = ? AND history_question_id IN ? ", user.ID, questionIDS).Error; err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "failed to remove old analytics for questions")
		return
	}

	var questionsAnalytics []models.HistoryQuestionAnalytic

	for _, questionAnalytic := range model.QuestionAnalytics {
		questionsAnalytics = append(questionsAnalytics, models.HistoryQuestionAnalytic{
			ID:                uuid.New(),
			HistoryQuestionID: questionAnalytic.HistoryQuestionID,
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
}

func GetWeeklyQuestionAnalytics(w http.ResponseWriter, r *http.Request) {
	user := utils.GetCurrentUser(r)
	if user == nil {
		utils.RespondWithError(w, http.StatusForbidden, "wrong token")
		return
	}

	//currentDay := time.Now().Weekday()

	// Calculate the number of days to subtract to get to the start of the week
	//daysToSubtract := int(currentDay)

	// Use Gorm to execute the SQL query
	var analytics []struct {
		DayOfWeek string `gorm:"day_of_week"`
		Count     int    `gorm:"count"`
	}

	if err := models.DB.Raw(`
        SELECT 
    date_trunc('day', created_at) AS day_of_week, 
    COUNT(*) AS count
FROM history_question_analytics
WHERE 
    user_id = ? AND
    created_at >= current_date - interval '7 days'
GROUP BY day_of_week
ORDER BY day_of_week;
    `, user.ID).Scan(&analytics).Error; err != nil {
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
