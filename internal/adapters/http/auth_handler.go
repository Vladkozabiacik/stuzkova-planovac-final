package http

import (
	"encoding/json"
	"net/http"
	"stuzkova-planovac/internal/core/dto"
	"stuzkova-planovac/internal/core/entity"
	"stuzkova-planovac/internal/core/service"

	"golang.org/x/crypto/bcrypt"
)

type AuthHandler struct {
	AuthService *service.AuthService
}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var registerDTO dto.RegisterDTO

	if err := json.NewDecoder(r.Body).Decode(&registerDTO); err != nil {
		http.Error(w, "Invalid input: "+err.Error(), http.StatusBadRequest)
		return
	}

	hashedPassword, err := hashPassword(registerDTO.Password)
	if err != nil {
		http.Error(w, "Error hashing password: "+err.Error(), http.StatusInternalServerError)
		return
	}

	user := &entity.User{
		Username: registerDTO.Username,
		Password: hashedPassword,
		Role:     registerDTO.Role,
	}

	if err := h.AuthService.Register(user); err != nil {
		http.Error(w, "Registration failed: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("User registered successfully"))
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var loginDTO dto.LoginDTO

	if err := json.NewDecoder(r.Body).Decode(&loginDTO); err != nil {
		http.Error(w, "Invalid input: "+err.Error(), http.StatusBadRequest)
		return
	}

	token, err := h.AuthService.Login(loginDTO.Username, loginDTO.Password)
	if err != nil {
		http.Error(w, "Login failed: "+err.Error(), http.StatusUnauthorized)
		return
	}

	response := map[string]string{
		"token": token,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func hashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}
