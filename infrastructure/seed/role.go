package seed

import (
	"case-management/internal/domain/model"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type RoleMap map[string]uuid.UUID

func SeedRole(db *gorm.DB) RoleMap {
	roleMap := make(RoleMap)

	roles := []model.Role{
		{Name: "Admin"},
		{Name: "Staff"},
		{Name: "Supervisor"},
		{Name: "AsstManager Up"},
	}

	for _, role := range roles {
		if err := db.FirstOrCreate(&role, model.Role{Name: role.Name}).Error; err != nil {
			panic(err)
		}
		roleMap[role.Name] = role.ID
	}
	return roleMap
}
