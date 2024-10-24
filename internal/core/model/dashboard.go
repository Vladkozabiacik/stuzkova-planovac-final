package model

import (
	"html/template"
	"net/http"

	"github.com/gorilla/mux"
)

type Dashboard struct {
	Role           string
	UserName       string
	WelcomeMessage string
}

// GetDashboardData returns data for the dashboard page
func GetDashboardData(role string, username string) Dashboard {
	var welcomeMessage string

	switch role {
	case "student":
		welcomeMessage = "Welcome to the Student Dashboard!"
	case "teacher":
		welcomeMessage = "Welcome to the Teacher Dashboard!"
	case "parent":
		welcomeMessage = "Welcome to the Parent Dashboard!"
	case "guest":
		welcomeMessage = "Welcome to the Guest Dashboard!"
	default:
		welcomeMessage = "Welcome!"
	}

	return Dashboard{
		Role:           role,
		UserName:       username,
		WelcomeMessage: welcomeMessage,
	}
}

// DashboardHandler renders the dashboard.html template with dynamic data
func DashboardHandler(w http.ResponseWriter, r *http.Request) {
	// Retrieve the role and username from the context (you may need to modify this based on how you store user info)
	role := r.Context().Value("role").(string)
	username := r.Context().Value("username").(string)

	// Get the dashboard data based on role and username
	dashboardData := GetDashboardData(role, username)

	// Parse the dashboard.html template
	tmpl, err := template.ParseFiles("templates/dashboard.html")
	if err != nil {
		http.Error(w, "Error loading template", http.StatusInternalServerError)
		return
	}

	// Execute the template and inject data
	err = tmpl.Execute(w, dashboardData)
	if err != nil {
		http.Error(w, "Error rendering template", http.StatusInternalServerError)
	}
}

// RegisterDashboardRoute registers the dashboard route
func RegisterDashboardRoute(router *mux.Router) {
	router.HandleFunc("/dashboard", DashboardHandler).Methods("GET")
}
