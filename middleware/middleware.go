package middleware

import (
	"github.com/gin-gonic/gin"
)

func GetUserID(c *gin.Context) uint {
	if userID, exists := c.Get("user_id"); exists {
		if id, ok := userID.(uint); ok {
			return id
		}
	}
	return 0
}

func GetUserRole(c *gin.Context) string {
	if role, exists := c.Get("role"); exists {
		if r, ok := role.(string); ok {
			return r
		}
	}
	return ""
}
