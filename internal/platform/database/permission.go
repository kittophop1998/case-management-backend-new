package database

import (
	"case-management/internal/domain/model"
	"fmt"

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

func (p *PermissionPg) UpdatePermission(
	ctx *gin.Context,
	departmentId uuid.UUID,
	sectionId uuid.UUID,
	reqs []model.UpdatePermissionRequest,
) error {
	return p.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		for _, req := range reqs {
			var permission model.Permission

			// หา permission หรือสร้างใหม่
			if err := tx.FirstOrCreate(&permission, model.Permission{Key: req.Permission}).Error; err != nil {
				return fmt.Errorf("permission key %s does not exist: %w", req.Permission, err)
			}

			// ดึง roles ที่ตรงกับชื่อ
			var roles []model.Role
			if err := tx.Where("name IN ?", req.Roles).Find(&roles).Error; err != nil {
				return err
			}
			if len(roles) == 0 {
				return fmt.Errorf("no valid roles found for permission %s", req.Permission)
			}

			// สร้าง RolePermission ถ้ายังไม่มี
			for _, role := range roles {
				rp := model.RolePermission{
					RoleID:       role.ID,
					PermissionID: permission.ID,
					SectionID:    sectionId,
					DepartmentID: departmentId,
				}

				if err := tx.FirstOrCreate(
					&rp,
					model.RolePermission{
						RoleID:       role.ID,
						PermissionID: permission.ID,
						SectionID:    sectionId,
						DepartmentID: departmentId,
					},
				).Error; err != nil {
					return err
				}
			}
		}
		return nil
	})
}

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
