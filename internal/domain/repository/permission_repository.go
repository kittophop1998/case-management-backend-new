package repository

import (
	"case-management/internal/domain/model"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type PermissionRepository interface {
	// GetAllPermissions(ctx *gin.Context, limit, offset int) ([]model.PermissionWithRolesResponse, error)
	GetAllPermissions(ctx *gin.Context, limit, offset int, permissionName string, sectionID, departmentID *uuid.UUID) ([]model.PermissionWithRolesResponse, error)
	UpdatePermission(
		ctx *gin.Context,
		req model.UpdatePermissionRequest,
		departmentMap map[string]uuid.UUID, // roleName -> departmentID
		sectionMap map[string]uuid.UUID, // roleName -> sectionID
	) error
	// CountPermissions(ctx *gin.Context) (int, error)
	CountPermissions(ctx *gin.Context, permissionName string, sectionID, departmentID *uuid.UUID) (int, error)
}
