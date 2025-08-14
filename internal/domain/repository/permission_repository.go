package repository

import (
	"case-management/internal/domain/model"

	"github.com/gin-gonic/gin"
)

type PermissionRepository interface {
	GetAllPermissions(ctx *gin.Context, limit, offset int) ([]model.PermissionWithRolesResponse, error)
	UpdatePermission(ctx *gin.Context, req model.UpdatePermissionRequest) error
	CountPermissions(ctx *gin.Context) (int, error)
}
