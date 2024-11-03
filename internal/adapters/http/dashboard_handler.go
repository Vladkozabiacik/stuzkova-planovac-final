package http

import (
	"log"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt"
)

func Dashboard(w http.ResponseWriter, r *http.Request) {
	authHeader := r.Header.Get("Authorization")
	// log.Println(authHeader)

	if !strings.HasPrefix(authHeader, "Bearer ") {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		log.Println("Authorization header missing or doesn't start with Bearer")
		return
	}

	tokenString := authHeader[7:]

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte("your_jwt_secret"), nil
	})

	if err != nil || !token.Valid {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		log.Printf("Failed to parse or invalid token: %v\n", err)
		return
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		role := claims["role"].(string)

		// log.Printf("Processing role: %s\n", role)

		// add template for dashboard and inject it with data based on user + his roles
		// important stuff, GDPR data accessibility, ...
		switch role {
		case "student":
			log.Println("Student dashboard accessed")
			w.Write([]byte("Welcome Student!"))
		case "teacher":
			log.Println("Teacher dashboard accessed")
			w.Write([]byte("Welcome Teacher!"))
		case "parent":
			log.Println("Parent dashboard accessed")
			w.Write([]byte("Welcome Parent!"))
		case "guest":
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
