package models

import (
	"github.com/google/uuid"
	"time"
)

type LessonAnalytic struct {
	ID                uuid.UUID `json:"id" gorm:"primary_key"`
	HistoryLessonID   uuid.UUID `json:"lesson_id"`
	UserID            uuid.UUID `json:"user_id"`
	RightAnswersCount int       `json:"right_answers_count"`
	QuestionsCount    int       `json:"questions_count"`
	TimeSpent         int       `json:"time_spent"`
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`
}

type QuestionAnalytic struct {
	ID            uuid.UUID `json:"id" gorm:"primary_key"`
	QuestionID    uuid.UUID `json:"question_id"`
	UserID        uuid.UUID `json:"user_id"`
	AnsweredRight bool      `json:"answered_right"`
	TimeSpent     int       `json:"time_spent"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

type Complaint struct {
	ID            uuid.UUID `json:"id" gorm:"primary_key"`
	UserID        uuid.UUID `json:"user_id"`
	LessonType    string    `json:"lesson_type"`
	LessonID      uuid.UUID `json:"lesson_id"`
	QuestionID    uuid.UUID `json:"question_id"`
	ComplaintText string    `json:"complaint_text"`
	IsSolved      bool      `json:"is_solved"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
