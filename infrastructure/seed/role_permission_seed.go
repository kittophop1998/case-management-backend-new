package seed

import (
	"case-management/internal/domain/model"
	"fmt"

	"gorm.io/gorm"
)

var rolePermissions = map[string][]string{
	"Staff": {
		"user.login",
		"user.profile",
		"user.assess",
		"user.customersearch",
		"user.verifycustomer",
		"user.customerdashboard",
	},
	"AsstManager Up": {
		"user.login",
		"user.profile",
		"user.assess",
		"user.customersearch",
		"user.verifycustomer",
		"user.customerdashboard",
		"case.management",
	},
	"Supervisor": {
		"user.login",
		"user.profile",
		"user.assess",
		"user.customersearch",
		"user.verifycustomer",
		"user.customerdashboard",
		"case.management",
	},
	"Admin": {
		"user.login",
		"user.profile",
		"user.assess",
		"user.customersearch",
		"user.verifycustomer",
		"user.customerdashboard",
		"case.management",
		"case.exporthistorical",
		"case.standardreport",
	},
}

func SeedRolePermission(db *gorm.DB, roleMap RoleMap, permissionMap PermissionMap) error {
	return db.Transaction(func(tx *gorm.DB) error {
		for roleName, permKeys := range rolePermissions {
			roleID, ok := roleMap[roleName]
			if !ok {
				return fmt.Errorf("role not found in roleMap: %s", roleName)
			}

			for _, permKey := range permKeys {
				permID, ok := permissionMap[permKey]
				if !ok {
					return fmt.Errorf("permission not found in permissionMap: %s", permKey)
				}

				rp := model.RolePermission{
					RoleID:       roleID,
					PermissionID: permID,
				}

				// ถ้ามีแล้วจะไม่ insert ซ้ำ
				if err := tx.FirstOrCreate(&rp, rp).Error; err != nil {
					return err
				}
			}
		}
		return nil
	})
}
