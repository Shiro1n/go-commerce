package repository

import (
	"database/sql"

	"github.com/shiro1n/go-commerce/user-service/internal/model"
)

// UserRepository interface
type UserRepository interface {
	FindAll() ([]model.User, error)
}

// userRepository struct
type userRepository struct {
	DB *sql.DB
}

// NewUserRepository creates a new UserRepository
func NewUserRepository(db *sql.DB) UserRepository {
	return &userRepository{DB: db}
}

// FindAll returns all users
func (r *userRepository) FindAll() ([]model.User, error) {
	rows, err := r.DB.Query("SELECT id, username, email FROM users")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []model.User
	for rows.Next() {
		var user model.User
		if err := rows.Scan(&user.ID, &user.Username, &user.Email); err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	return users, nil
}
