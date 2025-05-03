package utils

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type API struct {
	Success bool        `json:"success"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}

func APIResponse(ctx *gin.Context, message string, data interface{}) {
	ctx.JSON(http.StatusOK, API{
		Success: true,
		Message: message,
		Data:    data,
	})
}

func APIResponseWithError(ctx *gin.Context, err error) {
	if err == nil {
		ctx.JSON(http.StatusOK, API{
			Success: true,
			Message: "Operation successful",
		})
		return
	}

	ctx.JSON(http.StatusBadRequest, API{
		Success: false,
		Error:   err.Error(),
	})
}
