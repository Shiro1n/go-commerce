package repository

import (
	"database/sql"

	"github.com/shiro1n/go-commerce/user-service/internal/model"
)

type UserRepository interface {
	FindAll() ([]model.User, error)
	FindByUsername(username string) (*model.User, error)
	FindById(id int) (*model.User, error)
	Create(user *model.User) error
}

type userRepository struct {
	DB *sql.DB
}

func NewUserRepository(db *sql.DB) UserRepository {
	return &userRepository{DB: db}
}

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

func (r *userRepository) FindByUsername(username string) (*model.User, error) {
	row := r.DB.QueryRow("SELECT id, username, email, password FROM users WHERE username = $1", username)
	var user model.User
	if err := row.Scan(&user.ID, &user.Username, &user.Email, &user.Password); err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) FindById(userId int) (*model.User, error) {
	row := r.DB.QueryRow("SELECT id, username, email, password FROM users WHERE id = $1", userId)
	var user model.User
	if err := row.Scan(&user.ID, &user.Username, &user.Email, &user.Password); err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) Create(user *model.User) error {
	_, err := r.DB.Exec("INSERT INTO users (username, email, password) VALUES ($1, $2, $3)", user.Username, user.Email, user.Password)
	return err
}
