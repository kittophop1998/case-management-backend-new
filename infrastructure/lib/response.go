package lib

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Response[T any] struct {
	Data T `json:"data"`
}

type PaginatedResponse struct {
	Data       interface{} `json:"data"`
	Page       int         `json:"page"`
	Limit      int         `json:"limit"`
	Total      int         `json:"total"`
	TotalPages int         `json:"totalPages"`
}

func HandleResponse[T any](ctx *gin.Context, statusCode int, data T) {
	ctx.JSON(statusCode, Response[T]{Data: data})
}

// HandlePaginatedResponse ส่ง response แบบมี pagination
func HandlePaginatedResponse(ctx *gin.Context, data interface{}, page, limit, total int) {
	totalPages := (total + limit - 1) / limit
	ctx.JSON(http.StatusOK, PaginatedResponse{
		Data:       data,
		Page:       page,
		Limit:      limit,
		Total:      total,
		TotalPages: totalPages,
	})
}
