package seed

import (
	"case-management/internal/domain/model"

	"gorm.io/gorm"
)

func SeedCaseStatus(db *gorm.DB) {
	caseStatus := []model.CaseStatus{
		{Name: "created", Description: "-"},
		{Name: "progress", Description: "-"},
		{Name: "resolved", Description: "-"},
		{Name: "waiting", Description: "-"},
		{Name: "approved", Description: "-"},
		{Name: "closed", Description: "-"},
	}

	for _, status := range caseStatus {
		if err := db.FirstOrCreate(&status, model.CaseStatus{Name: status.Name}).Error; err != nil {
			panic(err)
		}
	}
}
