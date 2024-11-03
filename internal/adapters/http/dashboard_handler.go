// dashboard_handler.go
package http

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt"
)

type DashboardResponse struct {
	User     UserData     `json:"user"`
	Metadata MetadataInfo `json:"metadata"`
}

type UserData struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Role     string `json:"role"`
	Bio      string `json:"bio"`
}

type MetadataInfo struct {
	WelcomeMessage string   `json:"welcomeMessage"`
	Permissions    []string `json:"permissions"`
}

// DashboardData serves the JSON data for the dashboard.
func DashboardData(w http.ResponseWriter, r *http.Request) {
	// Add security headers
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.Header().Set("X-Frame-Options", "DENY")

	authHeader := r.Header.Get("Authorization")
	if !strings.HasPrefix(authHeader, "Bearer ") {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	tokenString := strings.TrimPrefix(authHeader, "Bearer ")
	claims := jwt.MapClaims{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte("your_jwt_secret"), nil
	})

	if err != nil || !token.Valid {
		http.Error(w, "Invalid token", http.StatusUnauthorized)
		return
	}

	username, _ := claims["username"].(string)

	// Construct response
	response := DashboardResponse{
		User: UserData{
			Username: username,
		},
	}

	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Printf("Error encoding response: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}

func getWelcomeMessage(role string) string {
	messages := map[string]string{
		"student": "Welcome to the Student Dashboard!",
		"teacher": "Welcome to the Teacher Dashboard!",
		"parent":  "Welcome to the Parent Dashboard!",
		"guest":   "Welcome to the Guest Dashboard!",
	}

	if msg, ok := messages[role]; ok {
		return msg
	}
	return "Welcome!"
}

func getRolePermissions(role string) []string {
	permissions := map[string][]string{
		"student": {"view_profile", "edit_profile", "view_occasions"},
		"teacher": {"view_profile", "edit_profile", "manage_occasions", "view_students"},
		"parent":  {"view_profile", "view_occasions"},
		"guest":   {"view_profile"},
	}

	if perms, ok := permissions[role]; ok {
		return perms
	}
	return []string{"view_profile"}
}
