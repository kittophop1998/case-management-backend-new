package seed

import (
	"case-management/internal/domain/model"

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

func SeedRolePermission(db *gorm.DB) error {
	return db.Transaction(func(tx *gorm.DB) error {
		for roleName, permKeys := range rolePermissions {
			var role model.Role
			if err := tx.Where("name = ?", roleName).First(&role).Error; err != nil {
				return err
			}

			var permissions []model.Permission
			for _, permKey := range permKeys {
				var permission model.Permission
				if err := tx.Where("key = ?", permKey).First(&permission).Error; err != nil {
					return err
				}
				permissions = append(permissions, permission)
			}

			// Replace ทั้งก้อนในรอบเดียว
			if err := tx.Model(&role).Association("Permissions").Replace(permissions); err != nil {
				return err
			}
		}
		return nil
	})
}
