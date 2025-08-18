package seed

import (
	"case-management/internal/domain/model"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type PermissionMap map[string]uuid.UUID

func SeedPermission(db *gorm.DB) PermissionMap {
	permissionMap := make(PermissionMap)

	permissions := []model.Permission{
		{Key: "user.login", Name: "Login"},
		{Key: "user.logout", Name: "Logout"},
		{Key: "user.management", Name: "User Management"},
		{Key: "user.profile", Name: "View profile"},
		{Key: "user.assess", Name: "Assess Control"},
		{Key: "user.customersearch", Name: "Customer Search"},
		{Key: "user.verifycustomer", Name: "Verify Customer"},
		{Key: "user.customerdashboard", Name: "Customer Dashboard"},
		{Key: "case.management", Name: "Case Management"},
		{Key: "case.exporthistorical", Name: "Export Case Historical"},
		{Key: "case.standardreport", Name: "Standard Report"},
	}

	for _, perm := range permissions {
		if err := db.FirstOrCreate(&perm, model.Permission{Key: perm.Key}).Error; err != nil {
			panic("failed to seed permissions: " + err.Error())
		}
		permissionMap[perm.Key] = perm.ID
	}

	return permissionMap
}
