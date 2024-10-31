package http

import (
	"encoding/json"
	"log"
	"net/http"
	"stuzkova-planovac/internal/core/dto"
	"stuzkova-planovac/internal/core/entity"
	"stuzkova-planovac/internal/core/service"
)

type AuthHandler struct {
	AuthService *service.AuthService
}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var registerDTO dto.RegisterDTO
	log.Printf("Registration attempt received")

	if err := json.NewDecoder(r.Body).Decode(&registerDTO); err != nil {
		log.Printf("Error decoding registration request: %v", err)
		http.Error(w, "Invalid input: "+err.Error(), http.StatusBadRequest)
		return
	}

	// Remove password hashing from here - let service handle it
	user := &entity.User{
		Username: registerDTO.Username,
		Password: registerDTO.Password, // Pass the plain password to service
		Role:     registerDTO.Role,
	}

	if err := h.AuthService.Register(user); err != nil {
		log.Printf("Error registering user: %v", err)
		http.Error(w, "Registration failed: "+err.Error(), http.StatusInternalServerError)
		return
	}

	log.Printf("User registered successfully: %s", registerDTO.Username)
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message": "User registered successfully"})
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var loginDTO dto.LoginDTO
	log.Printf("Login attempt received")

	if err := json.NewDecoder(r.Body).Decode(&loginDTO); err != nil {
		log.Printf("Error decoding login request: %v", err)
		http.Error(w, "Invalid input: "+err.Error(), http.StatusBadRequest)
		return
	}

	// Add debug logging
	log.Printf("Attempting login for username: %s", loginDTO.Username)

	token, err := h.AuthService.Login(loginDTO.Username, loginDTO.Password)
	if err != nil {
		log.Printf("Login failed for user %s: %v", loginDTO.Username, err)
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	response := map[string]string{
		"token": token,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
	log.Printf("Login successful for user: %s", loginDTO.Username)
}
