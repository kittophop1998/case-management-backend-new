package database

import (
	"case-management/internal/domain/model"

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

func (p *PermissionPg) GetAllPermissions(ctx *gin.Context, limit, offset int, permissionName string, sectionID, departmentID *uuid.UUID) ([]model.PermissionWithRolesResponse, int, error) {
	permQuery := p.db.WithContext(ctx).Model(&model.Permission{})
	if permissionName != "" {
		permQuery = permQuery.Where("name ILIKE ?", "%"+permissionName+"%")
	}

	var total int64
	if err := permQuery.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	var permissions []model.Permission
	if err := permQuery.
		Limit(limit).
		Offset(offset).
		Order("name ASC").
		Find(&permissions).Error; err != nil {
		return nil, 0, err
	}

	results := make([]model.PermissionWithRolesResponse, 0, len(permissions))
	permissionIDs := make([]uuid.UUID, 0, len(permissions))
	for _, perm := range permissions {
		results = append(results, model.PermissionWithRolesResponse{
			Permission: perm.Key,
			Name:       perm.Name,
			Roles:      []string{},
		})
		permissionIDs = append(permissionIDs, perm.ID)
	}

	if len(permissionIDs) > 0 {
		type roleRow struct {
			PermissionID uuid.UUID
			RoleName     string
		}

		roleQuery := p.db.WithContext(ctx).
			Table("role_permissions AS rp").
			Select("rp.permission_id, r.name AS role_name").
			Joins("JOIN roles AS r ON r.id = rp.role_id").
			Where("rp.permission_id IN ?", permissionIDs)

		if departmentID != nil {
			roleQuery = roleQuery.Where("rp.department_id = ?", *departmentID)
		}
		if sectionID != nil {
			roleQuery = roleQuery.Where("rp.section_id = ?", *sectionID)
		}

		var roleRows []roleRow
		if err := roleQuery.Find(&roleRows).Error; err != nil {
			return nil, 0, err
		}

		// เติม Roles ให้ตรงกับ Permission
		for i := range results {
			for _, row := range roleRows {
				if row.PermissionID == permissionIDs[i] {
					results[i].Roles = append(results[i].Roles, row.RoleName)
				}
			}
		}
	}

	return results, int(total), nil
}

func (p *PermissionPg) UpdatePermission(ctx *gin.Context, departmentId uuid.UUID, sectionId uuid.UUID, reqs []model.UpdatePermissionRequest) error {
	return p.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		permissionKeys := make([]string, 0, len(reqs))
		for _, req := range reqs {
			permissionKeys = append(permissionKeys, req.Permission)
		}

		var permissions []model.Permission
		if err := tx.Where("key IN ?", permissionKeys).Find(&permissions).Error; err != nil {
			return err
		}

		permMap := make(map[string]model.Permission)
		for _, perm := range permissions {
			permMap[perm.Key] = perm
		}

		for _, key := range permissionKeys {
			if _, ok := permMap[key]; !ok {
				perm := model.Permission{Key: key}
				if err := tx.Create(&perm).Error; err != nil {
					return err
				}
				permMap[key] = perm
			}
		}

		roleSet := make(map[string]struct{})
		for _, req := range reqs {
			for _, roleName := range req.Roles {
				roleSet[roleName] = struct{}{}
			}
		}

		roleNames := make([]string, 0, len(roleSet))
		for name := range roleSet {
			roleNames = append(roleNames, name)
		}

		var roles []model.Role
		if err := tx.Where("name IN ?", roleNames).Find(&roles).Error; err != nil {
			return err
		}

		roleMap := make(map[string]uuid.UUID)
		for _, r := range roles {
			roleMap[r.Name] = r.ID
		}

		newMappings := make([]model.RolePermission, 0)
		permissionIDs := make([]uuid.UUID, 0, len(reqs))

		for _, req := range reqs {
			perm := permMap[req.Permission]
			permissionIDs = append(permissionIDs, perm.ID)

			for _, roleName := range req.Roles {
				roleId, ok := roleMap[roleName]
				if !ok {
					continue
				}
				newMappings = append(newMappings, model.RolePermission{
					RoleID:       roleId,
					PermissionID: perm.ID,
					SectionID:    sectionId,
					DepartmentID: departmentId,
				})
			}
		}

		if len(permissionIDs) > 0 {
			if err := tx.Where("department_id = ? AND section_id = ? AND permission_id IN ?", departmentId, sectionId, permissionIDs).Delete(&model.RolePermission{}).Error; err != nil {
				return err
			}
		}

		if len(newMappings) > 0 {
			if err := tx.Create(&newMappings).Error; err != nil {
				return err
			}
		}

		return nil
	})
}
