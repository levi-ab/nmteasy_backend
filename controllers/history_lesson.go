package controllers

import (
	"encoding/json"
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

	json.NewEncoder(w).Encode(historyQuestions)
}

func GetHistoryLessons(w http.ResponseWriter, r *http.Request) {
	var historyLessons []models.HistoryLesson

	if err := models.DB.Order("created_at, title").Find(&historyLessons).Error; err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Failed to get lessons")
		return
	}

	groupedLessonsMap := make(map[string][]models.HistoryLessonWithProperTitle)

	for _, historyLesson := range historyLessons {
		properTitle, generalTitle := utils.FormatLessonTopic(historyLesson.Title)

		if lessons, ok := groupedLessonsMap[generalTitle]; ok {
			groupedLessonsMap[generalTitle] = append(lessons, models.HistoryLessonWithProperTitle{
				Title:       historyLesson.Title,
				ID:          historyLesson.ID,
				ProperTitle: properTitle,
				CreatedAt:   historyLesson.CreatedAt,
			})
		} else {
			groupedLessonsMap[generalTitle] = []models.HistoryLessonWithProperTitle{{
				Title:       historyLesson.Title,
				ID:          historyLesson.ID,
				ProperTitle: properTitle,
				CreatedAt:   historyLesson.CreatedAt,
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

	json.NewEncoder(w).Encode(result)
}
