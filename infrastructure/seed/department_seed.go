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
		{Key: "SYSTEM", Name: "System"},
		{Key: "CMS", Name: "Customer Service"},
		{Key: "Fraud", Name: "Engineering"},
		{Key: "Collection", Name: "Collection"},
		{Key: "EDP", Name: "EDP"},
		{Key: "AUTH", Name: "Authorize"},
		{Key: "Custodian", Name: "Custodian"},
		{Key: "Branch", Name: "Branch"},
	}

	for _, department := range departments {
		if err := db.FirstOrCreate(&department, model.Department{Name: department.Name}).Error; err != nil {
			panic(err)
		}
		departmentMap[department.Key] = department.ID
	}

	return departmentMap
}
