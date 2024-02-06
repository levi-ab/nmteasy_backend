package models

import (
	"github.com/google/uuid"
	"time"
)

type HistoryLessonAnalytic struct {
	ID                uuid.UUID `json:"id" gorm:"primary_key"`
	HistoryLessonID   uuid.UUID `json:"history_lesson_id"`
	UserID            uuid.UUID `json:"user_id"`
	RightAnswersCount int       `json:"right_answers_count"`
	QuestionsCount    int       `json:"questions_count"`
	TimeSpent         int       `json:"time_spent"`
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`

	HistoryLesson HistoryLesson `json:"-" gorm:"-"`
	User          User          `json:"-" gorm:"-"`
}
