package middleware

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt"
)

// CORS zlo
func CORSMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Set specific origin instead of wildcard
		origin := r.Header.Get("Origin")
		if origin == "" {
			origin = "http://localhost:8080"
		}

		w.Header().Set("Access-Control-Allow-Origin", origin)
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Authorization, Content-Type, Accept")
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		w.Header().Set("Access-Control-Max-Age", "3600")

		// Handle preflight requests
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

// debugging purposes
func LoggerMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Received request: %s %s", r.Method, r.URL.Path)

		for name, values := range r.Header {
			for _, value := range values {
				log.Printf("%s: %s", name, value)
			}
		}

		var bodyBytes []byte
		if r.Body != nil {
			bodyBytes, _ = io.ReadAll(r.Body)
			r.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
		}
		log.Printf("Body: %s", string(bodyBytes))

		next.ServeHTTP(w, r)
	})
}

// jwt token verification
func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var tokenString string

		// Log all headers for debugging
		/*
			log.Println("Request headers:")
			for name, values := range r.Header {
				log.Printf("%s: %v", name, values)
			}
		*/

		// Check Authorization header
		authHeader := r.Header.Get("Authorization")
		if authHeader != "" && strings.HasPrefix(authHeader, "Bearer ") {
			tokenString = authHeader[7:]
			// log.Println("Token found in Authorization header:", tokenString)
		} else {
			// Check cookie
			cookie, err := r.Cookie("jwtToken")
			if err == nil {
				tokenString = cookie.Value
				// log.Println("Token found in cookie:", tokenString)
				r.Header.Set("Authorization", "Bearer "+tokenString)
			} else {
				log.Printf("Cookie error: %v", err)
				http.Error(w, "Unauthorized - No token provided", http.StatusUnauthorized)
				return
			}
		}

		// Validate JWT token
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				log.Printf("Unexpected signing method: %v", token.Header["alg"])
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return []byte("your_jwt_secret"), nil
		})

		if err != nil {
			log.Printf("Token validation error: %v", err)
			http.Error(w, "Unauthorized - Invalid token", http.StatusUnauthorized)
			return
		}

		if !token.Valid {
			log.Println("Token is invalid")
			http.Error(w, "Unauthorized - Invalid token", http.StatusUnauthorized)
			return
		}

		// Token is valid, proceed
		// log.Println("Token validation successful")
		next.ServeHTTP(w, r)
	})
}
