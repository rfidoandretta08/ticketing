package controller

import (
	"net/http"
	"ticketing/dto"
	"ticketing/model"
	"ticketing/service"

	"github.com/gin-gonic/gin"
)

type AuthController struct {
	authService service.AuthService
}

func NewAuthController(authService service.AuthService) *AuthController {
	return &AuthController{authService: authService}
}

func (ac *AuthController) Register(c *gin.Context) {
	var req dto.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Konversi string ke model.Role
	var role model.Role
	if req.Role != "" {
		role = model.Role(req.Role) // Mengonversi string ke model.Role
	} else {
		// Jika tidak ada role, gunakan default 'user'
		role = model.Users
	}

	// Membuat objek User dari DTO
	user := &model.User{
		Name:     req.Name,
		Email:    req.Email,
		Password: req.Password, // Jangan lupa untuk hash password jika diperlukan
		Role:     role,         // Gunakan tipe model.Role
	}

	createdUser, err := ac.authService.Register(user)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Registration successful",
		"user":    createdUser,
	})
}

func (ac *AuthController) Login(c *gin.Context) {
	var request struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	token, user, err := ac.authService.Login(request.Email, request.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Login successful",
		"token":   token,
		"user":    user,
	})
}
