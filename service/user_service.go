package service

import (
	"ticketing/model"
	"ticketing/repository"
)

type UserService interface {
	GetUserByID(id uint) (*model.User, error)
	GetUserByEmail(email string) (*model.User, error)
	CreateUser(user *model.User) error
}

type userService struct {
	userRepo repository.UserRepository
}

func NewUserService(userRepo repository.UserRepository) UserService {
	return &userService{userRepo: userRepo}
}

func (s *userService) GetUserByID(id uint) (*model.User, error) {
	return s.userRepo.FindByID(id)
}

func (s *userService) GetUserByEmail(email string) (*model.User, error) {
	return s.userRepo.FindByEmail(email)
}

func (s *userService) CreateUser(user *model.User) error {
	return s.userRepo.Create(user)
}
