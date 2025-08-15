package seed

import (
	"case-management/internal/domain/model"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type DepartmentMap map[string]uuid.UUID

func SeedDepartment(db *gorm.DB) DepartmentMap {
	departmentMap := make(DepartmentMap)

	departments := []model.Department{
		{Name: "System"},
		{Name: "Finance"},
		{Name: "Engineering"},
		{Name: "Sales"},
		{Name: "Marketing"},
	}

	for _, department := range departments {
		if err := db.FirstOrCreate(&department, model.Department{Name: department.Name}).Error; err != nil {
			panic(err)
		}
		departmentMap[department.Name] = department.ID
	}

	return departmentMap
}
