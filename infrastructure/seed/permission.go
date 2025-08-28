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
		{Key: "view.user", Name: "View User Management"},
		{Key: "add.user", Name: "Add User (individual and import)"},
		{Key: "edit.user", Name: "Edit User"},
		{Key: "view.accesscontrol", Name: "View Access Control List"},
		{Key: "edit.accesscontrol", Name: "Edit Access Control"},
		{Key: "search.customer", Name: "Search For Customer"},
		{Key: "add.custnote", Name: "Add Customer Note"},
		{Key: "view.custnote", Name: "View Customer Note (list and detail)"},
		{Key: "add.case", Name: "Add New Case"},
		{Key: "view.case", Name: "View Case (list and detail)"},
		{Key: "edit.case", Name: "Edit Case"},
		{Key: "add.casenote", Name: "Add Case Note"},
		{Key: "view.inquirylog", Name: "View Inquiry Log"},
		{Key: "view.report", Name: "View Report"},
		{Key: "view.setting", Name: "View Setting"},
		{Key: "view.queue", Name: "View Queue Management"},
	}

	for _, perm := range permissions {
		if err := db.FirstOrCreate(&perm, model.Permission{Key: perm.Key}).Error; err != nil {
			panic("failed to seed permissions: " + err.Error())
		}
		permissionMap[perm.Key] = perm.ID
	}

	return permissionMap
}
