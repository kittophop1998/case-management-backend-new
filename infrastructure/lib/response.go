package lib

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// BaseResponse คือโครงสร้างพื้นฐานสำหรับทุก response
type BaseResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message,omitempty"`
}

// Response ใช้ส่งข้อมูลทั่วไป
type Response[T any] struct {
	BaseResponse
	Data T `json:"data"`
}

// PaginatedResponse ใช้ส่งข้อมูลแบบแบ่งหน้า (Pagination)
type PaginatedResponse[T any] struct {
	BaseResponse
	Data       []T `json:"data"`
	Page       int `json:"page"`
	Limit      int `json:"limit"`
	Total      int `json:"total"`
	TotalPages int `json:"totalPages"`
}

// HandleResponse ส่ง response ทั่วไป
func HandleResponse[T any](ctx *gin.Context, statusCode int, data T) {
	ctx.JSON(statusCode, Response[T]{
		BaseResponse: BaseResponse{
			Success: statusCode >= 200 && statusCode < 300,
		},
		Data: data,
	})
}

// HandlePaginatedResponse ส่ง response แบบมี pagination
func HandlePaginatedResponse[T any](ctx *gin.Context, page, limit, total int, data []T) {
	totalPages := 0
	if limit > 0 {
		totalPages = (total + limit - 1) / limit
	}

	ctx.JSON(http.StatusOK, PaginatedResponse[T]{
		BaseResponse: BaseResponse{
			Success: true,
		},
		Data:       data,
		Page:       page,
		Limit:      limit,
		Total:      total,
		TotalPages: totalPages,
	})
}
