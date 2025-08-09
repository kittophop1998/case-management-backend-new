package database

import (
	"case-management/internal/domain/model"
	"errors"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type PermissionPg struct {
	db *gorm.DB
}

func NewPermissionPg(db *gorm.DB) *PermissionPg {
	return &PermissionPg{db: db}
}

func (p *PermissionPg) GetAllPermissions(ctx *gin.Context, limit, offset int) ([]model.PermissionWithRolesResponse, error) {
	var permissions []model.Permission

	if err := p.db.WithContext(ctx).
		Preload("Roles").
		Limit(limit).
		Offset(offset).
		Order("name").
		Find(&permissions).Error; err != nil {
		return nil, err
	}

	var result []model.PermissionWithRolesResponse
	for _, p := range permissions {
		var roleNames []string
		for _, role := range p.Roles {
			roleNames = append(roleNames, role.Name)
		}

		if len(p.Roles) == 0 {
			roleNames = []string{}
		}

		result = append(result, model.PermissionWithRolesResponse{
			Permission: p.Key,
			Name:       p.Name,
			Roles:      roleNames,
		})
	}

	return result, nil
}

func (p *PermissionPg) UpdatePermission(ctx *gin.Context, req model.UpdatePermissionRequest) error {
	var permission model.Permission

	// Check if the permission exists
	if err := p.db.WithContext(ctx).Where("key = ?", req.Permission).First(&permission).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("permission key does not exist")
		}
	}

	var roles []model.Role

	// Check if roles exist
	if err := p.db.WithContext(ctx).Where("name IN ?", req.Roles).Find(&roles).Error; err != nil {
		return err
	}

	// If no roles are found, return an error
	if len(roles) == 0 {
		return errors.New("no valid roles found")
	}

	// Update the permission with the new roles
	if err := p.db.Model(&permission).Association("Roles").Replace(roles); err != nil {
		return err
	}

	return nil
}
