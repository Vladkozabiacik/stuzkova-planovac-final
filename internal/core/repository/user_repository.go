package repository

import (
	"database/sql"
	"fmt"
	"log"
	"stuzkova-planovac/internal/core/entity"
)

type UserRepository interface {
	Create(user *entity.User) error
	GetByUsername(username string) (*entity.User, error)
}

type userRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) Create(user *entity.User) error {
	log.Printf("Attempting to create user: %s", user.Username)

	// First check if user already exists
	existingUser := &entity.User{}
	err := r.db.QueryRow("SELECT username FROM users WHERE username = $1", user.Username).
		Scan(&existingUser.Username)
	if err != sql.ErrNoRows {
		if err == nil {
			log.Printf("User %s already exists", user.Username)
			return fmt.Errorf("username already exists")
		}
		// Only return error if it's not a "no rows" error
		log.Printf("Error checking existing user: %v", err)
		return fmt.Errorf("error checking existing user: %v", err)
	}

	// Set default role if empty
	if user.Role == "" {
		user.Role = "guest"
	}

	// Use QueryRow instead of Exec to get the inserted ID
	query := `
        INSERT INTO users (username, password, role) 
        VALUES ($1, $2, $3) 
        RETURNING id`

	err = r.db.QueryRow(query, user.Username, user.Password, user.Role).Scan(&user.ID)
	if err != nil {
		log.Printf("Error creating user: %v", err)
		return fmt.Errorf("error creating user: %v", err)
	}

	log.Printf("Successfully created user: %s with ID: %d", user.Username, user.ID)

	return nil
}

func (r *userRepository) GetByUsername(username string) (*entity.User, error) {
	log.Printf("Attempting to get user by username: %s", username)

	user := &entity.User{}
	query := `
        SELECT id, username, password, role 
        FROM users 
        WHERE username = $1`

	err := r.db.QueryRow(query, username).
		Scan(&user.ID, &user.Username, &user.Password, &user.Role)

	if err != nil {
		if err == sql.ErrNoRows {
			log.Printf("No user found with username: %s", username)
			return nil, fmt.Errorf("user not found")
		}
		log.Printf("Error getting user: %v", err)
		return nil, fmt.Errorf("error getting user: %v", err)
	}

	return user, nil
}
func (r *userRepository) DebugGetUserPassword(username string) {
	var storedHash string
	err := r.db.QueryRow("SELECT password FROM users WHERE username = $1", username).Scan(&storedHash)
	if err != nil {
		log.Printf("Debug - Error getting password hash: %v", err)
		return
	}
	log.Printf("Debug - Stored password hash for %s: %s", username, storedHash)
}
