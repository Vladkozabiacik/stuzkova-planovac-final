package http

import (
	"database/sql"
	"net/http"
	"stuzkova-planovac/internal/adapters/middleware"
	"stuzkova-planovac/internal/core/repository"
	"stuzkova-planovac/internal/core/service"

	"github.com/gorilla/mux"
)

func RegisterRoutes(router *mux.Router, db *sql.DB) {
	// Initialize auth service with user repository
	authService := service.NewAuthService(repository.NewUserRepository(db))
	authHandler := &AuthHandler{AuthService: authService}

	// Register API routes with logging and authentication middleware
	router.Handle("/dashboard", middleware.AuthMiddleware(http.HandlerFunc(Dashboard))).Methods("GET")

	// Public routes for registration and login
	router.HandleFunc("/register", authHandler.Register).Methods("POST")
	router.HandleFunc("/login", authHandler.Login).Methods("POST")

	// Serve index.html on root path
	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "static/templates/index.html")
	}).Methods("GET")

	// Serve static files from the /static directory
	router.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("static/"))))

	// Optionally, log the routes for verification
	// You can add this at the end for debugging purposes
	http.Handle("/", router)
}
