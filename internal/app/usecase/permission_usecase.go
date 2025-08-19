package usecase

import (
	"case-management/internal/domain/model"
	"case-management/internal/domain/repository"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type PermissionUseCase struct {
	permissionRepo repository.PermissionRepository
}

func NewPermissionUseCase(permissionRepo repository.PermissionRepository) *PermissionUseCase {
	return &PermissionUseCase{permissionRepo: permissionRepo}
}

func (p *PermissionUseCase) GetAllPermissions(ctx *gin.Context, page, limit int, permissionName string, sectionID, departmentID *uuid.UUID) ([]model.PermissionWithRolesResponse, int, error) {
	offset := (page - 1) * limit
	return p.permissionRepo.GetAllPermissions(ctx, limit, offset, permissionName, sectionID, departmentID)
}

func (p *PermissionUseCase) UpdatePermission(ctx *gin.Context, reqs []model.UpdatePermissionRequest, departmentId uuid.UUID, sectionId uuid.UUID) error {
	return p.permissionRepo.UpdatePermission(ctx, departmentId, sectionId, reqs)
}
