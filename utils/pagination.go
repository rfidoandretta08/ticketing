package utils

import (
	"math"
	"strconv"
	"ticketing/dto"

	"github.com/gin-gonic/gin"
)

func GeneratePagination(page, limit int, totalItems int64) dto.Pagination {
	totalPages := int64(0)
	if limit > 0 {
		totalPages = int64(math.Ceil(float64(totalItems) / float64(limit)))
	}

	return dto.Pagination{
		Page:       page,
		Limit:      limit,
		TotalItems: totalItems,
		TotalPages: totalPages,
	}
}

func ParsePaginationQuery(ctx *gin.Context) (int, int) {
	page, err := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	if err != nil || page < 1 {
		page = 1
	}

	limit, err := strconv.Atoi(ctx.DefaultQuery("limit", "10"))
	if err != nil || limit < 1 {
		limit = 10
	}

	return page, limit
}
