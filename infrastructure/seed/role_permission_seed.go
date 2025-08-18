package seed

import (
	"case-management/internal/domain/model"
	"fmt"

	"gorm.io/gorm"
)

var rolesConfig = map[string][]string{
	"System": {
		"view.profile",
		"view.user",
		"add.user",
		"edit.user",
		"view.accesscontrol",
		"edit.accesscontrol",
		"search.customer",
		"add.custnote",
		"view.custnote",
		"add.case",
		"view.case",
		"edit.case",
		"add.casenote",
	},
}

func SeedRolePermission(
	db *gorm.DB,
	roleMap RoleMap,
	permissionMap PermissionMap,
	departmentMap DepartmentMap, // roleName -> departmentID
	sectionMap SectionMap, // roleName -> sectionID
) error {
	return db.Transaction(func(tx *gorm.DB) error {

		for roleName, perms := range rolesConfig {
			roleID, ok := roleMap[roleName]
			if !ok {
				return fmt.Errorf("role not found in roleMap: %s", roleName)
			}

			deptID, ok := departmentMap[roleName]
			if !ok {
				return fmt.Errorf("department not found in departmentMap for role: %s", roleName)
			}

			secID, ok := sectionMap[roleName]
			if !ok {
				return fmt.Errorf("section not found in sectionMap for role: %s", roleName)
			}

			for _, permKey := range perms {
				permID, ok := permissionMap[permKey]
				if !ok {
					return fmt.Errorf("permission not found in permissionMap: %s", permKey)
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
		}
		return nil
	})
}
