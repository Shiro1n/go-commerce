package service

import (
	"errors"

	"github.com/shiro1n/go-commerce/user-service/internal/model"
	"github.com/shiro1n/go-commerce/user-service/internal/repository"
)

type UserService interface {
	GetAllUsers() ([]model.User, error)
	GetUserByUsername(username string) (*model.User, error)
	RegisterUser(username, email, password string) (*model.User, error)
	GetUserById(id int) (*model.User, error)
}

type userService struct {
	UserRepo repository.UserRepository
}

func NewUserService(userRepo repository.UserRepository) UserService {
	return &userService{UserRepo: userRepo}
}

func (s *userService) GetAllUsers() ([]model.User, error) {
	return s.UserRepo.FindAll()
}

func (s *userService) GetUserByUsername(username string) (*model.User, error) {
	return s.UserRepo.FindByUsername(username)
}

func (s *userService) GetUserById(userId int) (*model.User, error) {
	return s.UserRepo.FindById(userId)
}

func (s *userService) RegisterUser(username, email, password string) (*model.User, error) {
	// Check if the username already exists
	existingUser, _ := s.UserRepo.FindByUsername(username)
	if existingUser != nil {
		return nil, errors.New("username already taken")
	}

	// Create the new user
	user := &model.User{
		Username: username,
		Email:    email,
		Password: password, // Add proper password hashing
	}
	err := s.UserRepo.Create(user)
	if err != nil {
		return nil, err
	}

	return user, nil
}
