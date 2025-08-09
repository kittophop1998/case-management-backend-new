package usecase

import (
	"case-management/internal/domain/model"
	"case-management/internal/domain/repository"

	"github.com/gin-gonic/gin"
)

type PermissionUseCase struct {
	permissionRepo repository.PermissionRepository
}

func NewPermissionUseCase(permissionRepo repository.PermissionRepository) *PermissionUseCase {
	return &PermissionUseCase{permissionRepo: permissionRepo}
}

func (p *PermissionUseCase) GetAllPermissions(ctx *gin.Context, page, limit int) ([]model.PermissionWithRolesResponse, error) {
	offset := (page - 1) * limit

	permissions, err := p.permissionRepo.GetAllPermissions(ctx, limit, offset)
	if err != nil {
		return nil, err
	}

	return permissions, nil
}

func (p *PermissionUseCase) UpdatePermission(ctx *gin.Context, input model.UpdatePermissionRequest) error {
	return p.permissionRepo.UpdatePermission(ctx, input)
}
