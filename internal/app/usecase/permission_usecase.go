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

// func (p *PermissionUseCase) GetAllPermissions(ctx *gin.Context, page, limit int) ([]model.PermissionWithRolesResponse, int, error) {
// 	offset := (page - 1) * limit

// 	permissions, err := p.permissionRepo.GetAllPermissions(ctx, limit, offset)
// 	if err != nil {
// 		return nil, 0, err
// 	}

// 	total, err := p.permissionRepo.CountPermissions(ctx)
// 	if err != nil {
// 		return nil, 0, err
// 	}

// 	return permissions, total, nil
// }

func (p *PermissionUseCase) GetAllPermissions(ctx *gin.Context, page, limit int, permissionName string, sectionID, departmentID *uuid.UUID) ([]model.PermissionWithRolesResponse, int, error) {
	offset := (page - 1) * limit

	permissions, err := p.permissionRepo.GetAllPermissions(ctx, limit, offset, permissionName, sectionID, departmentID)
	if err != nil {
		return nil, 0, err
	}

	total, err := p.permissionRepo.CountPermissions(ctx, permissionName, sectionID, departmentID)
	if err != nil {
		return nil, 0, err
	}

	return permissions, total, nil
}

func (p *PermissionUseCase) UpdatePermission(
	ctx *gin.Context,
	reqs []model.UpdatePermissionRequest,
	departmentId uuid.UUID,
	sectionId uuid.UUID,
) error {
	return p.permissionRepo.UpdatePermission(ctx, departmentId, sectionId, reqs)
}
