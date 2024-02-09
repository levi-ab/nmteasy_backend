package migrated_models

import (
	"github.com/google/uuid"
	"time"
)

type HistoryQuestion struct {
	ID              uuid.UUID `json:"id" gorm:"primary_key"`
	QuestionText    string    `json:"question_text"`
	QuestionImage   string    `json:"question_image"`
	Answers         string    `json:"answers"`
	Type            string    `json:"type"`
	RightAnswer     string    `json:"right_answer"`
	Topic           string    `json:"topic"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
	HistoryLessonID uuid.UUID `json:"lesson_id"`
}

type HistoryLesson struct {
	ID        uuid.UUID         `json:"id" gorm:"primary_key"`
	Title     string            `json:"title"`
	Questions []HistoryQuestion `json:"questions"`
	TimeSpent int               `json:"time_spent"`
	CreatedAt time.Time         `json:"created_at"`
	UpdatedAt time.Time         `json:"updated_at"`
}

type HistoryQuestionExplanation struct {
	ID                uuid.UUID `json:"id" gorm:"primary_key"`
	Explanation       string    `json:"explanation"`
	HistoryQuestionID uuid.UUID `json:"question_id"`
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`

	HistoryQuestion HistoryQuestion `json:"question"`
}

type UkrainianQuestion struct {
	ID                uuid.UUID `json:"id" gorm:"primary_key"`
	QuestionText      string    `json:"question_text"`
	QuestionImage     string    `json:"question_image"`
	Answers           string    `json:"answers"`
	Type              string    `json:"type"`
	RightAnswer       string    `json:"right_answer"`
	Topic             string    `json:"topic"`
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`
	UkrainianLessonID uuid.UUID `json:"ukrainian_lesson_id"`
}

type UkrainianLesson struct {
	ID        uuid.UUID           `json:"id" gorm:"primary_key"`
	Title     string              `json:"title"`
	Questions []UkrainianQuestion `json:"questions"`
	TimeSpent int                 `json:"time_spent"`
	CreatedAt time.Time           `json:"created_at"`
	UpdatedAt time.Time           `json:"updated_at"`
}

type UkrainianQuestionExplanation struct {
	ID                  uuid.UUID `json:"id" gorm:"primary_key"`
	Explanation         string    `json:"explanation"`
	UkrainianQuestionID uuid.UUID `json:"ukrainian_question_id"`
	CreatedAt           time.Time `json:"created_at"`
	UpdatedAt           time.Time `json:"updated_at"`

	UkrainianQuestion UkrainianQuestion `json:"ukrainian_question"`
}

type BiologyQuestion struct {
	ID              uuid.UUID `json:"id" gorm:"primary_key"`
	QuestionText    string    `json:"question_text"`
	QuestionImage   string    `json:"question_image"`
	Answers         string    `json:"answers"`
	Type            string    `json:"type"`
	RightAnswer     string    `json:"right_answer"`
	Topic           string    `json:"topic"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
	BiologyLessonID uuid.UUID `json:"biology_lesson_id"`
}

type BiologyLesson struct {
	ID        uuid.UUID         `json:"id" gorm:"primary_key"`
	Title     string            `json:"title"`
	Questions []BiologyQuestion `json:"questions"`
	TimeSpent int               `json:"time_spent"`
	CreatedAt time.Time         `json:"created_at"`
	UpdatedAt time.Time         `json:"updated_at"`
}

type BiologyQuestionExplanation struct {
	ID                uuid.UUID `json:"id" gorm:"primary_key"`
	Explanation       string    `json:"explanation"`
	BiologyQuestionID uuid.UUID `json:"biology_question_id"`
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`

	BiologyQuestion BiologyQuestion `json:"biology_question"`
}
