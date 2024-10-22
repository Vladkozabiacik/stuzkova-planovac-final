package model

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
)

// Dashboard represents the data structure for the dashboard.
type Dashboard struct {
	Role           string `json:"role"`
	WelcomeMessage string `json:"welcome_message"`
}

// GetDashboardData generates dashboard data based on the user's role.
func GetDashboardData(role string) Dashboard {
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
		WelcomeMessage: welcomeMessage,
	}
}

// DashboardHandler handles requests to the dashboard endpoint.
func DashboardHandler(w http.ResponseWriter, r *http.Request) {
	// Retrieve the role from the context or request (assuming it's passed as a JWT claim or similar).
	role := r.Context().Value("role").(string) // Replace with your method of extracting user role

	// Get dashboard data based on user role
	dashboardData := GetDashboardData(role)

	// Respond with the dashboard data (in JSON format)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(dashboardData)
}

// RegisterDashboardRoute registers the dashboard route with the provided router.
func RegisterDashboardRoute(router *mux.Router) {
	router.HandleFunc("/dashboard", DashboardHandler).Methods("GET")
}
