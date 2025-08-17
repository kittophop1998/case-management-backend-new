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
	// 1️⃣ Query paginated permissions
	permQuery := p.db.WithContext(ctx).Model(&model.Permission{})
	if permissionName != "" {
		permQuery = permQuery.Where("name ILIKE ?", "%"+permissionName+"%")
	}

	var permissions []model.Permission
	if err := permQuery.
		Limit(limit).
		Offset(offset).
		Find(&permissions).Error; err != nil {
		return nil, err
	}

	// 2️⃣ Build permission map
	permMap := make(map[uuid.UUID]*model.PermissionWithRolesResponse)
	for _, perm := range permissions {
		permMap[perm.ID] = &model.PermissionWithRolesResponse{
			Permission: perm.Key,
			Name:       perm.Name,
			Roles:      []string{},
		}
	}

	// 3️⃣ Query role_permissions for these permissions
	type roleRow struct {
		PermissionID uuid.UUID
		RoleName     string
	}

	permissionIDs := make([]uuid.UUID, 0, len(permissions))
	for _, perm := range permissions {
		permissionIDs = append(permissionIDs, perm.ID)
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
		return nil, err
	}

	// 4️⃣ Append roles to permissions
	for _, row := range roleRows {
		if permResp, ok := permMap[row.PermissionID]; ok {
			permResp.Roles = append(permResp.Roles, row.RoleName)
		}
	}

	// 5️⃣ Convert map to slice
	results := make([]model.PermissionWithRolesResponse, 0, len(permMap))
	for _, v := range permMap {
		results = append(results, *v)
	}

	return results, nil
}

func (p *PermissionPg) UpdatePermission(ctx *gin.Context, departmentId uuid.UUID, sectionId uuid.UUID, reqs []model.UpdatePermissionRequest) error {
	return p.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// 1) หา/สร้าง permissions ทีเดียว
		permissionKeys := make([]string, 0, len(reqs))
		for _, req := range reqs {
			permissionKeys = append(permissionKeys, req.Permission)
		}

		var permissions []model.Permission
		if err := tx.Where("key IN ?", permissionKeys).Find(&permissions).Error; err != nil {
			return err
		}

		// map เก็บ key -> permission
		permMap := make(map[string]model.Permission)
		for _, perm := range permissions {
			permMap[perm.Key] = perm
		}

		// สร้างใหม่ถ้ายังไม่มี
		for _, key := range permissionKeys {
			if _, ok := permMap[key]; !ok {
				perm := model.Permission{Key: key}
				if err := tx.Create(&perm).Error; err != nil {
					return err
				}
				permMap[key] = perm
			}
		}

		// 2) หา roles ที่เกี่ยวข้องทั้งหมด
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

		// 3) เตรียม RolePermission ใหม่
		newMappings := make([]model.RolePermission, 0)
		permissionIDs := make([]uuid.UUID, 0, len(reqs))

		for _, req := range reqs {
			perm := permMap[req.Permission]
			permissionIDs = append(permissionIDs, perm.ID)

			for _, roleName := range req.Roles {
				roleId, ok := roleMap[roleName]
				if !ok {
					continue // skip role ที่หาไม่เจอ
				}
				newMappings = append(newMappings, model.RolePermission{
					RoleID:       roleId,
					PermissionID: perm.ID,
					SectionID:    sectionId,
					DepartmentID: departmentId,
				})
			}
		}

		// 4) ลบ mapping เดิมทั้งหมดของ permissionIds + section + department
		if len(permissionIDs) > 0 {
			if err := tx.Where("department_id = ? AND section_id = ? AND permission_id IN ?", departmentId, sectionId, permissionIDs).Delete(&model.RolePermission{}).Error; err != nil {
				return err
			}
		}

		fmt.Printf("Updating permissions for department %s, section %s with %d new mappings\n", departmentId, sectionId, len(newMappings))
		fmt.Printf("New mappings: %+v\n", newMappings)
		// 5) Insert ใหม่ทั้งหมด (bulk)
		if len(newMappings) > 0 {
			if err := tx.Create(&newMappings).Error; err != nil {
				return err
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
