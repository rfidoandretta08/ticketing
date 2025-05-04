package service

import (
	"errors"
	"fmt"
	"strings"
	"ticketing/model"
	"ticketing/repository"
	"ticketing/utils"

	"golang.org/x/crypto/bcrypt"
)

type AuthService interface {
	Login(email, password string) (string, *model.User, error)
	Register(user *model.User) (*model.User, error)
}

type authService struct {
	userRepo repository.UserRepository
}

func NewAuthService(userRepo repository.UserRepository) AuthService {
	return &authService{userRepo: userRepo}
}

func (s *authService) Login(email, password string) (string, *model.User, error) {
	user, err := s.userRepo.FindByEmail(email)
	if err != nil {
		return "", nil, errors.New("invalid email or password")
	}

	password = strings.TrimSpace(password)

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return "", nil, errors.New("invalid email or password")
	}

	token, err := utils.GenerateToken(user.ID, string(user.Role), user.Email)
	if err != nil {
		return "", nil, err
	}

	return token, user, nil
}

func (s *authService) Register(user *model.User) (*model.User, error) {
	existing, _ := s.userRepo.FindByEmail(user.Email)
	if existing != nil && existing.ID != 0 {
		return nil, errors.New("email already registered")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	user.Password = string(hashedPassword)

	fmt.Println("DEBUG: Hashed password (to save):", user.Password)

	if err := s.userRepo.Create(user); err != nil {
		return nil, err
	}
	return user, nil
}
