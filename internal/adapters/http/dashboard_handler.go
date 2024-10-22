package http

import (
	"log"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt"
)

func Dashboard(w http.ResponseWriter, r *http.Request) {
	// Get the Authorization header
	authHeader := r.Header.Get("Authorization")
	log.Println(authHeader)
	// Check if the header is provided and starts with "Bearer "
	if !strings.HasPrefix(authHeader, "Bearer ") {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		log.Println("Authorization header missing or doesn't start with Bearer")
		return
	}

	// Extract the token string
	tokenString := authHeader[7:] // The token starts right after "Bearer "

	// Parse the JWT token
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte("your_jwt_secret"), nil
	})

	if err != nil || !token.Valid {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		log.Printf("Failed to parse or invalid token: %v\n", err)
		return
	}

	// Extract claims
	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		role := claims["role"].(string)

		// Log the role
		log.Printf("Processing role: %s\n", role)

		switch role {
		case "student":
			// Handle student dashboard
			log.Println("Student dashboard accessed")
			w.Write([]byte("Welcome Student!"))
		case "teacher":
			// Handle teacher dashboard
			log.Println("Teacher dashboard accessed")
			w.Write([]byte("Welcome Teacher!"))
		case "parent":
			// Handle parent dashboard
			log.Println("Parent dashboard accessed")
			w.Write([]byte("Welcome Parent!"))
		case "guest":
			// Handle guest dashboard
			log.Println("Guest dashboard accessed")
			w.Write([]byte("Welcome Guest!"))
		default:
			log.Printf("Unauthorized role: %s\n", role)
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
		}
	} else {
		log.Println("Failed to extract claims from token")
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
	}
}
