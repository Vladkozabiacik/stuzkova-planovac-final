package main

import (
	"log"
	"net/http"

	"stuzkova-planovac/internal/adapters/db"
	httpAdapter "stuzkova-planovac/internal/adapters/http"
	"stuzkova-planovac/internal/config"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	cfg := config.LoadConfig()

	db, err := db.Connect(cfg.DatabaseURL)
	if err != nil {
		log.Fatal(err)
	}

	router := mux.NewRouter()
	httpAdapter.RegisterRoutes(router, db)
	log.Println("Server starting on :8080")
	if err := http.ListenAndServe(":8080", router); err != nil {
		log.Fatal(err)
	}
}
