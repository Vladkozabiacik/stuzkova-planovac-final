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
	// auth services
	authService := service.NewAuthService(repository.NewUserRepository(db))
	authHandler := &AuthHandler{AuthService: authService}

	// Public routes
	registerPublicRoutes(router, authHandler)

	// Protected routes
	registerProtectedRoutes(router)

	// Static files
	router.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("static/"))))

	http.Handle("/", router)
}

func registerPublicRoutes(router *mux.Router, authHandler *AuthHandler) {
	// Auth routes
	router.HandleFunc("/register", authHandler.Register).Methods("POST")
	router.HandleFunc("/login", authHandler.Login).Methods("POST")

	// Root path
	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "static/templates/index.html")
	}).Methods("GET")
}

func registerProtectedRoutes(router *mux.Router) {
	// Create a subrouter for protected routes
	protected := router.PathPrefix("").Subrouter()

	// Apply auth middleware to all routes in this subrouter
	protected.Use(middleware.AuthMiddleware)

	// Dashboard routes
	protected.HandleFunc("/dashboard", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "static/templates/dashboard.html")
	}).Methods("GET")
	protected.HandleFunc("/dashboard-data", Dashboard).Methods("GET")

	// Add more protected routes here
	protected.HandleFunc("/profile", HandleProfile).Methods("GET", "PUT")
	protected.HandleFunc("/settings", HandleSettings).Methods("GET", "PUT")
	protected.HandleFunc("/api/user-data", HandleUserData).Methods("GET")
}

// Handler functions for new protected routes
func HandleProfile(w http.ResponseWriter, r *http.Request) {
	// Profile handling logic
}

func HandleSettings(w http.ResponseWriter, r *http.Request) {
	// Settings handling logic
}

func HandleUserData(w http.ResponseWriter, r *http.Request) {
	// User data handling logic
}
