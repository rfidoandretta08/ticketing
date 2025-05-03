package controller

import (
	"net/http"
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
	var user model.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	createdUser, err := ac.authService.Register(&user)
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
