package service

import (
	"fmt"
	"log"
	"stuzkova-planovac/internal/core/entity"
	"stuzkova-planovac/internal/core/repository"
	"time"

	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	repo repository.UserRepository
}

func NewAuthService(repo repository.UserRepository) *AuthService {
	return &AuthService{repo: repo}
}

func (s *AuthService) Register(user *entity.User) error {
	log.Printf("Register - Original password: %s", user.Password)

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Printf("Error hashing password: %v", err)
		return err
	}

	user.Password = string(hashedPassword)
	return s.repo.Create(user)
}

func (s *AuthService) Login(username, password string) (string, error) {
	log.Printf("Attempting login for username: %s", username)

	user, err := s.repo.GetByUsername(username)
	if err != nil {
		log.Printf("Error retrieving user: %v", err)
		return "", fmt.Errorf("invalid credentials")
	}

	// Compare the provided password with the stored hash
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		log.Printf("Password comparison failed for user %s: %v", username, err)
		return "", fmt.Errorf("invalid credentials")
	}

	// Generate JWT token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":       user.ID,
		"username": user.Username,
		"exp":      time.Now().Add(time.Hour * 24).Unix(),
	})

	tokenString, err := token.SignedString([]byte("your_jwt_secret"))
	if err != nil {
		log.Printf("Error generating token: %v", err)
		return "", fmt.Errorf("error generating token")
	}

	log.Printf("Login successful for user: %s", username)
	return tokenString, nil
}
