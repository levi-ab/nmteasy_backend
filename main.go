package main

import (
	"github.com/joho/godotenv"
	"net/http"
	"nmteasy_backend/common"
	"nmteasy_backend/controllers"
	"nmteasy_backend/models"
	"os"
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

	server.ListenAndServe()
}
