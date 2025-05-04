package repository

import (
	"fmt"
	"ticketing/model"

	"gorm.io/gorm"
)

type UserRepository interface {
	Create(user *model.User) error
	FindByEmail(email string) (*model.User, error)
	FindByID(id uint) (*model.User, error)
	FindAllUsers(page, limit int) ([]model.User, int64, error)
	GetAllUsers() ([]model.User, error)
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) Create(user *model.User) error {
	return r.db.Create(user).Error
}

func (r *userRepository) FindByEmail(email string) (*model.User, error) {
	var user model.User
	fmt.Println("DEBUG: Looking for email:", email)
	err := r.db.Where("email = ?", email).First(&user).Error
	if err != nil {
		fmt.Println("DEBUG: Error from DB:", err)
		return nil, err
	}
	fmt.Println("DEBUG: Found user ID:", user.ID)
	return &user, nil
}

func (r *userRepository) FindByID(id uint) (*model.User, error) {
	var user model.User
	err := r.db.First(&user, id).Error
	return &user, err
}

func (r *userRepository) GetAllUsers() ([]model.User, error) {
	var users []model.User
	err := r.db.Find(&users).Error
	return users, err
}

func (r *userRepository) FindAllUsers(page, limit int) ([]model.User, int64, error) {
	var users []model.User
	var total int64

	// Menghitung jumlah total data
	if err := r.db.Model(&model.User{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Mengambil data dengan limit dan offset
	offset := (page - 1) * limit
	if err := r.db.Offset(offset).Limit(limit).Find(&users).Error; err != nil {
		return nil, 0, err
	}

	return users, total, nil
}
