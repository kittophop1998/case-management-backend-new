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
		{Key: "view.profile", Name: "View profile"},
		{Key: "view.user", Name: "View User List"},
		{Key: "add.user", Name: "add user (individual and import)"},
	}

	for _, perm := range permissions {
		if err := db.FirstOrCreate(&perm, model.Permission{Key: perm.Key}).Error; err != nil {
			panic("failed to seed permissions: " + err.Error())
		}
		permissionMap[perm.Key] = perm.ID
	}

	return permissionMap
}
