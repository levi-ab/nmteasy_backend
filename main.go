package main

import (
	"github.com/joho/godotenv"
	"net/http"
	"nmteasy_backend/common"
	"nmteasy_backend/controllers"
	"nmteasy_backend/models"
	"nmteasy_backend/utils"
	"os"
	"time"
)

func main() {
	godotenv.Load()

	handler := controllers.New()

	server := &http.Server{
		Addr:    "0.0.0.0:8008",
		Handler: handler,
	}

	models.ConnectDatabase()
	common.SECRET_KEY = os.Getenv("SECRET_KEY")

	utils.GenerateRandomUsers(models.DB)

	go func() {
		ticker := time.NewTicker(7 * 24 * time.Hour)

		// Run the function immediately before waiting for the first tick
		utils.ResetLeagues()

		// Loop to run the function every time the ticker ticks
		for {
			select {
			case <-ticker.C:
				utils.ResetLeagues()
			}
		}
	}()

	server.ListenAndServe()
}
