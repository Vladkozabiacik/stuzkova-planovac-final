package http

import (
	"encoding/json"
	"net/http"
	"stuzkova-planovac/internal/core/dto"
	"stuzkova-planovac/internal/core/entity"
	"stuzkova-planovac/internal/core/service"

	"golang.org/x/crypto/bcrypt" // Import bcrypt for password hashing
)

type AuthHandler struct {
	AuthService *service.AuthService
}

// Register handles user registration
func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var registerDTO dto.RegisterDTO

	// Decode the incoming JSON body
	if err := json.NewDecoder(r.Body).Decode(&registerDTO); err != nil {
		http.Error(w, "Invalid input: "+err.Error(), http.StatusBadRequest)
		return
	}

	// Hash the password for secure storage
	hashedPassword, err := hashPassword(registerDTO.Password)
	if err != nil {
		http.Error(w, "Error hashing password: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Create a new user entity
	user := &entity.User{
		Username: registerDTO.Username,
		Password: hashedPassword, // Store the hashed password
		Role:     registerDTO.Role,
	}

	// Attempt to register the user using the AuthService
	if err := h.AuthService.Register(user); err != nil {
		http.Error(w, "Registration failed: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Successfully created user
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("User registered successfully"))
}

// Login handles user authentication
func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var loginDTO dto.LoginDTO

	// Decode the incoming JSON body
	if err := json.NewDecoder(r.Body).Decode(&loginDTO); err != nil {
		http.Error(w, "Invalid input: "+err.Error(), http.StatusBadRequest)
		return
	}

	// Attempt to authenticate and retrieve the token
	token, err := h.AuthService.Login(loginDTO.Username, loginDTO.Password)
	if err != nil {
		http.Error(w, "Login failed: "+err.Error(), http.StatusUnauthorized)
		return
	}

	// Respond with the token in the body
	response := map[string]string{
		"token": token,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// hashPassword hashes the password using bcrypt
func hashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}
