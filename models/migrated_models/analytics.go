package migrated_models

import (
	"github.com/google/uuid"
	"time"
)

type HistoryLessonAnalytic struct {
	ID                uuid.UUID `json:"id" gorm:"primary_key"`
	HistoryLessonID   uuid.UUID `json:"lesson_id"`
	UserID            uuid.UUID `json:"user_id"`
	RightAnswersCount int       `json:"right_answers_count"`
	QuestionsCount    int       `json:"questions_count"`
	TimeSpent         int       `json:"time_spent"`
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`

	HistoryLesson HistoryLesson `json:"-" gorm:"-"`
	User          User          `json:"-" gorm:"-"`
}

type HistoryQuestionAnalytic struct {
	ID                uuid.UUID `json:"id" gorm:"primary_key"`
	HistoryQuestionID uuid.UUID `json:"question_id" gorm:"history_question_id"`
	UserID            uuid.UUID `json:"user_id"`
	AnsweredRight     bool      `json:"answered_right"`
	TimeSpent         int       `json:"time_spent"`
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`

	HistoryQuestion HistoryQuestion `json:"-" gorm:"-"`
	User            User            `json:"-" gorm:"-"`
}

type UkrainianLessonAnalytic struct {
	ID                uuid.UUID `json:"id" gorm:"primary_key"`
	UkrainianLessonID uuid.UUID `json:"lesson_id"`
	UserID            uuid.UUID `json:"user_id"`
	RightAnswersCount int       `json:"right_answers_count"`
	QuestionsCount    int       `json:"questions_count"`
	TimeSpent         int       `json:"time_spent"`
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`

	UkrainianLesson UkrainianLesson `json:"-" gorm:"-"`
	User            User            `json:"-" gorm:"-"`
}

type UkrainianQuestionAnalytic struct {
	ID                  uuid.UUID `json:"id" gorm:"primary_key"`
	UkrainianQuestionID uuid.UUID `json:"question_id" gorm:"ukrainian_question_id"`
	UserID              uuid.UUID `json:"user_id"`
	AnsweredRight       bool      `json:"answered_right"`
	TimeSpent           int       `json:"time_spent"`
	CreatedAt           time.Time `json:"created_at"`
	UpdatedAt           time.Time `json:"updated_at"`

	UkrainianQuestion UkrainianQuestion `json:"-" gorm:"-"`
	User              User              `json:"-" gorm:"-"`
}

type BiologyLessonAnalytic struct {
	ID                uuid.UUID `json:"id" gorm:"primary_key"`
	BiologyLessonID   uuid.UUID `json:"lesson_id"`
	UserID            uuid.UUID `json:"user_id"`
	RightAnswersCount int       `json:"right_answers_count"`
	QuestionsCount    int       `json:"questions_count"`
	TimeSpent         int       `json:"time_spent"`
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`

	BiologyLesson BiologyLesson `json:"-" gorm:"-"`
	User          User          `json:"-" gorm:"-"`
}

type BiologyQuestionAnalytic struct {
	ID                uuid.UUID `json:"id" gorm:"primary_key"`
	BiologyQuestionID uuid.UUID `json:"question_id" gorm:"biology_question_id"`
	UserID            uuid.UUID `json:"user_id"`
	AnsweredRight     bool      `json:"answered_right"`
	TimeSpent         int       `json:"time_spent"`
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`

	BiologyQuestion BiologyQuestion `json:"-" gorm:"-"`
	User            User            `json:"-" gorm:"-"`
}
