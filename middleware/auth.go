package middleware

import (
	"net/http"
	"strings"
	"ticketing/utils"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware(allowedRoles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Missing or invalid token"})
			c.Abort()
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		// Parse token dan ambil claims
		claims, err := utils.ParseToken(tokenString)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
			c.Abort()
			return
		}

		roleClaim, ok := claims["role"].(string)
		if !ok {
			c.JSON(http.StatusForbidden, gin.H{"error": "Invalid role format"})
			c.Abort()
			return
		}

		// Cek apakah role dari token termasuk dalam daftar allowedRoles
		isAllowed := false
		for _, role := range allowedRoles {
			if roleClaim == role {
				isAllowed = true
				break
			}
		}

		if !isAllowed {
			c.JSON(http.StatusForbidden, gin.H{"error": "Forbidden: Insufficient permissions"})
			c.Abort()
			return
		}

		// Set user_id dan role ke context
		if idFloat, ok := claims["id"].(float64); ok {
			c.Set("user_id", uint(idFloat))
		}
		c.Set("role", roleClaim)

		c.Next()
	}
}
