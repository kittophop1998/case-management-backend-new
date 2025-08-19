package seed

import (
	"case-management/internal/domain/model"
	"fmt"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type RolesConfigItem struct {
	RoleKey       string
	PermissionKey string
	DepartmentKey string
	SectionKey    string
}

var rolesConfig = []RolesConfigItem{
	{RoleKey: "System", PermissionKey: "view.user", DepartmentKey: "SYSTEM", SectionKey: "SYSTEM"},
	{RoleKey: "System", PermissionKey: "add.user", DepartmentKey: "SYSTEM", SectionKey: "SYSTEM"},
	{RoleKey: "System", PermissionKey: "edit.user", DepartmentKey: "SYSTEM", SectionKey: "SYSTEM"},
	{RoleKey: "System", PermissionKey: "view.accesscontrol", DepartmentKey: "SYSTEM", SectionKey: "SYSTEM"},
	{RoleKey: "System", PermissionKey: "edit.accesscontrol", DepartmentKey: "SYSTEM", SectionKey: "SYSTEM"},
	{RoleKey: "System", PermissionKey: "search.customer", DepartmentKey: "SYSTEM", SectionKey: "SYSTEM"},
	{RoleKey: "System", PermissionKey: "add.custnote", DepartmentKey: "SYSTEM", SectionKey: "SYSTEM"},
	{RoleKey: "System", PermissionKey: "view.custnote", DepartmentKey: "SYSTEM", SectionKey: "SYSTEM"},
	{RoleKey: "System", PermissionKey: "add.case", DepartmentKey: "SYSTEM", SectionKey: "SYSTEM"},
	{RoleKey: "System", PermissionKey: "view.case", DepartmentKey: "SYSTEM", SectionKey: "SYSTEM"},
	{RoleKey: "System", PermissionKey: "edit.case", DepartmentKey: "SYSTEM", SectionKey: "SYSTEM"},
	{RoleKey: "System", PermissionKey: "add.casenote", DepartmentKey: "SYSTEM", SectionKey: "SYSTEM"},
}

func SeedRolePermission(
	db *gorm.DB,
	roleMap map[string]uuid.UUID,
	permissionMap map[string]uuid.UUID,
	departmentMap map[string]uuid.UUID,
	sectionMap map[string]uuid.UUID,
) error {
	return db.Transaction(func(tx *gorm.DB) error {

		for _, cfg := range rolesConfig {
			roleID, ok := roleMap[cfg.RoleKey]
			if !ok {
				return fmt.Errorf("role not found in roleMap: %s", cfg.RoleKey)
			}

			permID, ok := permissionMap[cfg.PermissionKey]
			if !ok {
				return fmt.Errorf("permission not found in permissionMap: %s", cfg.PermissionKey)
			}

			deptID, ok := departmentMap[cfg.DepartmentKey]
			if !ok {
				return fmt.Errorf("department not found in departmentMap for role: %s", cfg.RoleKey)
			}

			secID, ok := sectionMap[cfg.SectionKey]
			if !ok {
				return fmt.Errorf("section not found in sectionMap for role: %s", cfg.RoleKey)
			}

			rp := model.RolePermission{
				RoleID:       roleID,
				PermissionID: permID,
				DepartmentID: deptID,
				SectionID:    secID,
			}

			condition := map[string]interface{}{
				"role_id":       roleID,
				"permission_id": permID,
				"department_id": deptID,
				"section_id":    secID,
			}

			if err := tx.FirstOrCreate(&rp, condition).Error; err != nil {
				return err
			}
		}

		return nil
	})
}
