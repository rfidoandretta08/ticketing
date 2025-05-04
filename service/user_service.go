package service

import (
	"ticketing/dto"
	"ticketing/model"
	"ticketing/repository"
	"ticketing/utils"
)

type UserService interface {
	GetUserByID(id uint) (*model.User, error)
	GetUserByEmail(email string) (*model.User, error)
	CreateUser(user *model.User) error
	GetAllUsers(page, limit int) ([]model.User, dto.Pagination, error)
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

func (s *userService) GetAllUsers(page, limit int) ([]model.User, dto.Pagination, error) {
	users, totalItems, err := s.userRepo.FindAllUsers(page, limit)
	if err != nil {
		return nil, dto.Pagination{}, err
	}

	pagination := utils.GeneratePagination(page, limit, totalItems)
	return users, pagination, nil
}
