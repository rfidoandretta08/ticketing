package middleware

import (
	"ticketing/utils"

	"github.com/gin-gonic/gin"
)

func ErrorHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.Next()

		if len(ctx.Errors) > 0 {
			err := ctx.Errors.Last()
			// Call APIResponse with just two arguments: context and message
			if err.Err != nil {
				utils.APIResponseWithError(ctx, err.Err)
			} else {
				// You can replace this with a message like "No errors occurred" if no error is present
				utils.APIResponse(ctx, "No errors occurred", nil)
			}
		}
	}
}
