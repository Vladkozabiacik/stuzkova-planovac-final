package middleware

import (
	"bytes"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt"
)

// CORSMiddleware sets the necessary headers to handle CORS
func CORSMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Set CORS headers
		w.Header().Set("Access-Control-Allow-Origin", "*")                                // Allow all origins, or replace "*" with your specific domain
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS") // Allowed HTTP methods
		w.Header().Set("Access-Control-Allow-Headers", "Authorization, Content-Type")     // Allowed headers

		// If it's an OPTIONS request, respond with 200 OK and return (preflight request)
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		// Call the next handler if it's not an OPTIONS request
		next.ServeHTTP(w, r)
	})
}

// LoggerMiddleware logs incoming requests
func LoggerMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Log method and URL
		log.Printf("Received request: %s %s", r.Method, r.URL.Path)

		// Log headers
		for name, values := range r.Header {
			for _, value := range values {
				log.Printf("%s: %s", name, value)
			}
		}

		// Log the request body
		var bodyBytes []byte
		if r.Body != nil {
			bodyBytes, _ = ioutil.ReadAll(r.Body)                 // Read the body
			r.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes)) // Restore the body
		}
		log.Printf("Body: %s", string(bodyBytes))

		next.ServeHTTP(w, r) // Call the next handler
	})
}

// AuthMiddleware checks for a valid JWT token in the request
func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var tokenString string

		// Check Authorization header first
		authHeader := r.Header.Get("Authorization")
		if authHeader != "" && strings.HasPrefix(authHeader, "Bearer ") {
			tokenString = authHeader[7:] // Extract token from header
		} else {
			// If Authorization header is missing or invalid, check cookies
			cookie, err := r.Cookie("jwtToken")
			if err == nil {
				tokenString = cookie.Value // Get token from the cookie
				log.Println("Token found in cookie:", tokenString)

				// Set the Authorization header for downstream handlers
				r.Header.Set("Authorization", "Bearer "+tokenString)
			} else {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				log.Println("Authorization header missing and cookie not found")
				return
			}
		}

		// Parse and validate the JWT token
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			// Ensure the token's signing method is as expected (HMAC in this case)
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, http.ErrAbortHandler
			}
			// Return the secret key to validate the token
			return []byte("your_jwt_secret"), nil
		})

		// If there was an error or the token is invalid, respond with Unauthorized
		if err != nil || !token.Valid {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			log.Printf("Invalid token: %v", err)
			return
		}

		// If token is valid, call the next handler
		log.Println("AuthMiddleware executed")
		next.ServeHTTP(w, r)
	})
}
