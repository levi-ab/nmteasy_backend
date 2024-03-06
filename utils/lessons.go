package utils

import (
	"fmt"
	"nmteasy_backend/models"
)

func GetRandomQuestionsByType(lessonType string, limit int) ([]models.Question, error) {
	var questions []models.Question

	rows, err := models.DB.
		Raw(fmt.Sprintf(`SELECT id, question_text, question_image, answers, type, right_answer, topic, created_at, updated_at, %s_lesson_id as lesson_id
    FROM %s_questions  ORDER BY RANDOM()  LIMIT %d`, lessonType, lessonType, limit)).Rows()

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var question models.Question
		if err := models.DB.ScanRows(rows, &question); err != nil {
			return nil, err
		}

		questions = append(questions, question)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return questions, nil
}
