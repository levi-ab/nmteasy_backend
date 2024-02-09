package models

import (
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"nmteasy_backend/models/migrated_models"
	"os"
)

var DB *gorm.DB

func ConnectDatabase() {
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_PORT"),
	)
	database, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		panic("Failed to connect to database")
	}

	database.AutoMigrate(&migrated_models.User{})
	database.AutoMigrate(&Complaint{})
	database.AutoMigrate(&migrated_models.League{})

	database.AutoMigrate(&migrated_models.HistoryQuestion{})
	database.AutoMigrate(&migrated_models.HistoryLesson{})
	database.AutoMigrate(&migrated_models.HistoryQuestionExplanation{})
	database.AutoMigrate(&migrated_models.HistoryLessonAnalytic{})
	database.AutoMigrate(&migrated_models.HistoryQuestionAnalytic{})

	database.AutoMigrate(&migrated_models.UkrainianQuestion{})
	database.AutoMigrate(&migrated_models.UkrainianLesson{})
	database.AutoMigrate(&migrated_models.UkrainianQuestionExplanation{})
	database.AutoMigrate(&migrated_models.UkrainianLessonAnalytic{})
	database.AutoMigrate(&migrated_models.UkrainianQuestionAnalytic{})

	database.AutoMigrate(&migrated_models.BiologyQuestion{})
	database.AutoMigrate(&migrated_models.BiologyLesson{})
	database.AutoMigrate(&migrated_models.BiologyQuestionExplanation{})
	database.AutoMigrate(&migrated_models.BiologyLessonAnalytic{})
	database.AutoMigrate(&migrated_models.BiologyQuestionAnalytic{})

	DB = database
}
