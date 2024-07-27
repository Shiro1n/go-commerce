package repository

import (
	"database/sql"

	"github.com/shiro1n/go-commerce/user-service/internal/model"
)

type UserRepository interface {
	FindById(userId int) (*model.User, error)
	Create(user *model.User) error
	Update(user *model.User) error
	Delete(userId int) error
}

type userRepository struct {
	DB *sql.DB
}

func NewUserRepository(db *sql.DB) UserRepository {
	return &userRepository{DB: db}
}

func (r *userRepository) FindById(userId int) (*model.User, error) {
	var user model.User
	err := r.DB.QueryRow("SELECT id, username, email FROM users WHERE id = $1", userId).Scan(&user.ID, &user.Username, &user.Email)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) Create(user *model.User) error {
	_, err := r.DB.Exec("INSERT INTO users (username, email, password) VALUES ($1, $2, $3)", user.Username, user.Email, user.Password)
	return err
}

func (r *userRepository) Update(user *model.User) error {
	_, err := r.DB.Exec("UPDATE users SET username = $1, email = $2, password = $3 WHERE id = $4", user.Username, user.Email, user.Password, user.ID)
	return err
}

func (r *userRepository) Delete(userId int) error {
	_, err := r.DB.Exec("DELETE FROM users WHERE id = $1", userId)
	return err
}
