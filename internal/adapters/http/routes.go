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
	// Apply CORS middleware to all routes
	router.Use(middleware.CORSMiddleware)

	// Apply logger middleware if needed
	// router.Use(middleware.LoggerMiddleware)

	// auth services
	authService := service.NewAuthService(repository.NewUserRepository(db))
	authHandler := &AuthHandler{AuthService: authService}

	// Public routes
	registerPublicRoutes(router, authHandler)

	// Protected routes
	registerProtectedRoutes(router)

	// Static files
	router.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("static/"))))

	// Important: Handle OPTIONS requests at the root level
	router.Methods("OPTIONS").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	http.Handle("/", router)
}

func registerPublicRoutes(router *mux.Router, authHandler *AuthHandler) {
	// Auth routes - ensure they handle OPTIONS
	router.HandleFunc("/register", authHandler.Register).Methods("POST", "OPTIONS")
	router.HandleFunc("/login", authHandler.Login).Methods("POST", "OPTIONS")

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

	// Dashboard route
	protected.HandleFunc("/dashboard", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "static/templates/dashboard.html")
	}).Methods("GET")
	// TODO: add functions for these routes
	// // Fetch dashboard data dynamically
	protected.HandleFunc("/dashboard-data", DashboardData).Methods("GET")

	// // Profile routes
	// protected.HandleFunc("/profile", HandleProfile).Methods("GET", "PUT", "OPTIONS")
	// protected.HandleFunc("/profile/edit", EditProfileForm).Methods("GET")
	// protected.HandleFunc("/profile/update", UpdateProfile).Methods("POST")

	// // Occasion routes
	// protected.HandleFunc("/occasions", ListOccasions).Methods("GET")
	// protected.HandleFunc("/occasions/new", NewOccasionForm).Methods("GET")
	// protected.HandleFunc("/occasions/create", CreateOccasion).Methods("POST")
	// protected.HandleFunc("/occasions/edit/{id}", EditOccasionForm).Methods("GET")
	// protected.HandleFunc("/occasions/update/{id}", UpdateOccasion).Methods("POST")
	// protected.HandleFunc("/occasions/delete/{id}", DeleteOccasion).Methods("DELETE")
}
