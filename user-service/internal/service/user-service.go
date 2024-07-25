package service

import (
	"github.com/shiro1n/go-commerce/user-service/internal/model"
	"github.com/shiro1n/go-commerce/user-service/internal/repository"
)

// UserService interface
type UserService interface {
	GetAllUsers() ([]model.User, error)
}

// userService struct
type userService struct {
	UserRepo repository.UserRepository
}

// NewUserService creates a new UserService
func NewUserService(userRepo repository.UserRepository) UserService {
	return &userService{UserRepo: userRepo}
}

// GetAllUsers returns all users
func (s *userService) GetAllUsers() ([]model.User, error) {
	return s.UserRepo.FindAll()
}
