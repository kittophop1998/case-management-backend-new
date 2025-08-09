package seed

import (
	"case-management/internal/domain/model"
	"log"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

func SeedPermission(db *gorm.DB) map[string]uuid.UUID {
	permissionMap := make(map[string]uuid.UUID)

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

	for _, p := range permissions {
		var permission model.Permission
		err := db.Where("key = ?", p.Key).FirstOrCreate(&permission, p).Error
		if err != nil {
			log.Printf("failed to seed permission %s: %v", p.Name, err)
		} else {
			permissionMap[p.Key] = permission.ID
		}
	}

	return permissionMap
}
