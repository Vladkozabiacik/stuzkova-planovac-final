package middleware

import (
	"bytes"
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt"
)

// CORS zlo
func CORSMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Set CORS headers
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Authorization, Content-Type")

		// preflight request
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

		authHeader := r.Header.Get("Authorization")
		if authHeader != "" && strings.HasPrefix(authHeader, "Bearer ") {
			tokenString = authHeader[7:]
		} else {
			cookie, err := r.Cookie("jwtToken")
			if err == nil {
				tokenString = cookie.Value

				r.Header.Set("Authorization", "Bearer "+tokenString)
			} else {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				log.Println("Authorization header missing and cookie not found")
				return
			}
		}

		// jwt validation
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			// is HMAC?
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, http.ErrAbortHandler
			}
			return []byte("your_jwt_secret"), nil
		})

		if err != nil || !token.Valid {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			log.Printf("Invalid token: %v", err)
			return
		}

		log.Println("AuthMiddleware executed")
		next.ServeHTTP(w, r)
	})
}
