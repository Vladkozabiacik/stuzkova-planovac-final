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

	// core
	router.Handle("/dashboard-data", middleware.AuthMiddleware(http.HandlerFunc(Dashboard))).Methods("GET")
	router.HandleFunc("/dashboard", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "static/templates/dashboard.html")
	}).Methods("GET")
	// auth
	router.HandleFunc("/register", authHandler.Register).Methods("POST")
	router.HandleFunc("/login", authHandler.Login).Methods("POST")

	// root path
	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "static/templates/index.html")
	}).Methods("GET")

	// static
	router.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("static/"))))

	http.Handle("/", router)
}
