package main

import (
	"github.com/joho/godotenv"
	"net/http"
	"nmteasy_backend/controllers"
	"nmteasy_backend/models"
)

func main() {
	godotenv.Load()

	handler := controllers.New()

	server := &http.Server{
		Addr:    "0.0.0.0:8008",
		Handler: handler,
	}

	models.ConnectDatabase()

	server.ListenAndServe()
}
