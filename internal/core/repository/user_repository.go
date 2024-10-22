package repository

import (
	"database/sql"
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
	if user.Role == "" {
		user.Role = "guest"
	}

	_, err := r.db.Exec("INSERT INTO users (username, password, role) VALUES ($1, $2, $3)",
		user.Username, user.Password, user.Role)
	return err
}

func (r *userRepository) GetByUsername(username string) (*entity.User, error) {
	user := &entity.User{}
	err := r.db.QueryRow("SELECT id, username, password, role FROM users WHERE username = $1", username).
		Scan(&user.ID, &user.Username, &user.Password, &user.Role)
	if err != nil {
		return nil, err
	}
	return user, nil
}
