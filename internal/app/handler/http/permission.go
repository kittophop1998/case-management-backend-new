package http

import (
	"case-management/infrastructure/lib"
	"case-management/internal/app/usecase"
	"case-management/internal/domain/model"
	"net/http"

	"github.com/gin-gonic/gin"
)

type PermissionHandler struct {
	UseCase usecase.PermissionUseCase
}

func (h *PermissionHandler) GetAllPermissions(ctx *gin.Context) {
	limit, err := getLimit(ctx)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid limit parameter"})
		return
	}

	page, err := getPage(ctx)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid page parameter"})
		return
	}

	permissions, total, err := h.UseCase.GetAllPermissions(ctx, page, limit)
	if err != nil {
		ctx.JSON(500, gin.H{"error": err.Error()})
		return
	}

	lib.HandlePaginatedResponse(ctx, page, limit, total, permissions)
}

func (h *PermissionHandler) UpdatePermission(ctx *gin.Context) {
	var req model.UpdatePermissionRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.UseCase.UpdatePermission(ctx, req); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Permission updated successfully"})
}
