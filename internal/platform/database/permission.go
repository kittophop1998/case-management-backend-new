package database

import (
	"case-management/internal/domain/model"
	"sort"

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

func (p *PermissionPg) GetAllPermissions(ctx *gin.Context, limit, offset int, permissionName string, sectionID, departmentID *uuid.UUID) ([]model.PermissionWithRolesResponse, int, int, error) {

	// --- query base สำหรับ permissions ทั้งหมดตาม filter ---
	permQuery := p.db.WithContext(ctx).Model(&model.Permission{})
	if permissionName != "" {
		permQuery = permQuery.Where("name ILIKE ?", "%"+permissionName+"%")
	}

	// --- นับ total permissions ทั้งหมด ---
	var total int64
	if err := permQuery.Count(&total).Error; err != nil {
		return nil, 0, 0, err
	}

	// --- ดึง permission ทั้งหมด ---
	var allPermissions []model.Permission
	if err := permQuery.Order("name ASC").Find(&allPermissions).Error; err != nil {
		return nil, 0, 0, err
	}

	// --- ดึง role สำหรับทุก permission ---
	permissionIDs := make([]uuid.UUID, len(allPermissions))
	for i, perm := range allPermissions {
		permissionIDs[i] = perm.ID
	}

	roleMap := make(map[uuid.UUID][]string)
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
			return nil, 0, 0, err
		}

		for _, row := range roleRows {
			roleMap[row.PermissionID] = append(roleMap[row.PermissionID], row.RoleName)
		}
	}

	// --- สร้าง results ---
	results := make([]model.PermissionWithRolesResponse, len(allPermissions))
	permWithRoles := make(map[uuid.UUID]struct{})
	for i, perm := range allPermissions {
		roles := roleMap[perm.ID]
		if roles == nil {
			roles = []string{} // ถ้าไม่มี role ให้เป็น slice ว่าง
		}

		results[i] = model.PermissionWithRolesResponse{
			Permission: perm.Key,
			Name:       perm.Name,
			Roles:      roles,
		}

		if len(roles) > 0 {
			permWithRoles[perm.ID] = struct{}{}
		}
	}

	permissionRoleCount := len(permWithRoles)

	// --- sort ตามจำนวน role จากมากไปน้อย แล้วตามชื่อ ---
	sort.SliceStable(results, func(i, j int) bool {
		if len(results[i].Roles) != len(results[j].Roles) {
			return len(results[i].Roles) > len(results[j].Roles)
		}
		return results[i].Name < results[j].Name
	})

	// --- slice สำหรับ pagination ---
	start := offset
	end := offset + limit
	if start > len(results) {
		start = len(results)
	}
	if end > len(results) {
		end = len(results)
	}
	paginated := results[start:end]

	return paginated, len(allPermissions), permissionRoleCount, nil
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
