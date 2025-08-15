package database

import (
	"case-management/internal/domain/model"
	"errors"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type PermissionPg struct {
	db *gorm.DB
}

func NewPermissionPg(db *gorm.DB) *PermissionPg {
	return &PermissionPg{db: db}
}

// func (p *PermissionPg) GetAllPermissions(ctx *gin.Context, limit, offset int) ([]model.PermissionWithRolesResponse, error) {
// 	var permissions []model.Permission

// 	if err := p.db.WithContext(ctx).
// 		Preload("Roles").
// 		Limit(limit).
// 		Offset(offset).
// 		Order("name").
// 		Find(&permissions).Error; err != nil {
// 		return nil, err
// 	}

// 	var result []model.PermissionWithRolesResponse
// 	for _, p := range permissions {
// 		var roleNames []string
// 		for _, role := range p.Roles {
// 			roleNames = append(roleNames, role.Name)
// 		}

// 		if len(p.Roles) == 0 {
// 			roleNames = []string{}
// 		}

// 		result = append(result, model.PermissionWithRolesResponse{
// 			Permission: p.Key,
// 			Name:       p.Name,
// 			Roles:      roleNames,
// 		})
// 	}

// 	return result, nil
// }

func (p *PermissionPg) GetAllPermissions(ctx *gin.Context, limit, offset int, permissionName string, sectionID, departmentID *uuid.UUID) ([]model.PermissionWithRolesResponse, error) {

	var permissions []model.Permission

	query := p.db.WithContext(ctx).
		Model(&model.Permission{}).
		Preload("Roles").
		Joins("LEFT JOIN role_permissions rp ON rp.permission_id = permissions.id")

	// กรองชื่อ permission ถ้าส่งมา
	if permissionName != "" {
		query = query.Where("permissions.name ILIKE ?", "%"+permissionName+"%")
	}

	// กรอง Section ถ้ามี
	if sectionID != nil {
		query = query.Where("rp.section_id = ?", *sectionID)
	}

	// กรอง Department ถ้ามี
	if departmentID != nil {
		query = query.Where("rp.department_id = ?", *departmentID)
	}

	if err := query.
		Limit(limit).
		Offset(offset).
		Group("permissions.id").
		Order("permissions.name").
		Find(&permissions).Error; err != nil {
		return nil, err
	}

	var result []model.PermissionWithRolesResponse
	for _, perm := range permissions {
		var roleNames []string
		for _, role := range perm.Roles {
			roleNames = append(roleNames, role.Name)
		}
		if len(perm.Roles) == 0 {
			roleNames = []string{}
		}
		result = append(result, model.PermissionWithRolesResponse{
			Permission: perm.Key,
			Name:       perm.Name,
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

// func (p *PermissionPg) CountPermissions(ctx *gin.Context) (int, error) {
// 	var count int64
// 	if err := p.db.WithContext(ctx).Model(&model.Permission{}).Count(&count).Error; err != nil {
// 		return 0, err
// 	}
// 	return int(count), nil
// }

func (p *PermissionPg) CountPermissions(ctx *gin.Context, permissionName string, sectionID, departmentID *uuid.UUID) (int, error) {
	var count int64
	query := p.db.WithContext(ctx).Model(&model.Permission{}).
		Joins("LEFT JOIN role_permissions rp ON rp.permission_id = permissions.id")

	if permissionName != "" {
		query = query.Where("permissions.name ILIKE ?", "%"+permissionName+"%")
	}
	if sectionID != nil {
		query = query.Where("rp.section_id = ?", *sectionID)
	}
	if departmentID != nil {
		query = query.Where("rp.department_id = ?", *departmentID)
	}

	if err := query.Count(&count).Error; err != nil {
		return 0, err
	}

	return int(count), nil
}
