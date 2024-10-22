package model

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
)

// create dashboard.html template display and inject it with data based on user
type Dashboard struct {
	Role           string `json:"role"`
	WelcomeMessage string `json:"welcome_message"`
}

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

func DashboardHandler(w http.ResponseWriter, r *http.Request) {
	role := r.Context().Value("role").(string)

	dashboardData := GetDashboardData(role)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(dashboardData)
}

func RegisterDashboardRoute(router *mux.Router) {
	router.HandleFunc("/dashboard", DashboardHandler).Methods("GET")
}
