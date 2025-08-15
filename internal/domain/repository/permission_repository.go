package repository

import (
	"case-management/internal/domain/model"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type PermissionRepository interface {
	GetAllPermissions(ctx *gin.Context, limit, offset int, permissionName string, sectionID, departmentID *uuid.UUID) ([]model.PermissionWithRolesResponse, error)
	UpdatePermission(ctx *gin.Context, departmentId uuid.UUID, sectionId uuid.UUID, reqs []model.UpdatePermissionRequest) error
	CountPermissions(ctx *gin.Context, permissionName string, sectionID, departmentID *uuid.UUID) (int, error)
}
