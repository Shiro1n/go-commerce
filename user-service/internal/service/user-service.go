package service

import (
	"github.com/shiro1n/go-commerce/user-service/internal/model"
	"github.com/shiro1n/go-commerce/user-service/internal/repository"
)

type UserService interface {
	GetUserById(userId int) (*model.User, error)
	CreateUser(user *model.User) error
	UpdateUser(user *model.User) error
	DeleteUser(userId int) error
}

type userService struct {
	UserRepo repository.UserRepository
}

func NewUserService(userRepo repository.UserRepository) UserService {
	return &userService{UserRepo: userRepo}
}

func (s *userService) GetUserById(userId int) (*model.User, error) {
	return s.UserRepo.FindById(userId)
}

func (s *userService) CreateUser(user *model.User) error {
	return s.UserRepo.Create(user)
}

func (s *userService) UpdateUser(user *model.User) error {
	return s.UserRepo.Update(user)
}

func (s *userService) DeleteUser(userId int) error {
	return s.UserRepo.Delete(userId)
}
