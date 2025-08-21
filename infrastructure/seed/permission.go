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
		{Key: "add.user", Name: "Add user (individual and import)"},
		{Key: "edit.user", Name: "Edit user"},
		{Key: "view.accesscontrol", Name: "View access control list"},
		{Key: "edit.accesscontrol", Name: "Edit access control"},
		{Key: "search.customer", Name: "Search For Customer"},
		{Key: "add.custnote", Name: "Add customer note"},
		{Key: "view.custnote", Name: "View customer note (list and detail)"},
		{Key: "add.case", Name: "Add new case"},
		{Key: "view.case", Name: "View case (list and detail)"},
		{Key: "edit.case", Name: "Edit case"},
		{Key: "add.casenote", Name: "Add case note"},
		{Key: "view.inquirylog", Name: "View inquiry log"},
		{Key: "view.report", Name: "View report"},
		{Key: "view.setting", Name: "View setting"},
	}

	for _, perm := range permissions {
		if err := db.FirstOrCreate(&perm, model.Permission{Key: perm.Key}).Error; err != nil {
			panic("failed to seed permissions: " + err.Error())
		}
		permissionMap[perm.Key] = perm.ID
	}

	return permissionMap
}
