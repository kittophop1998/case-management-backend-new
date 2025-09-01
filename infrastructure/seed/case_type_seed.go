package seed

import (
	"case-management/internal/domain/model"

	"gorm.io/gorm"
)

func SeedCaseType(db *gorm.DB) {
	caseTypes := []model.CaseTypes{
		{Name: "Inquiry and disposition", Group: "Inquiry", Description: "-"},
		{Name: "Change Address", Group: "Change Info", Description: "-"},
		{Name: "Change Mobile no.", Group: "Change Info", Description: "-"},
		{Name: "Change Passport", Group: "Change Info", Description: "-"},
		{Name: "Change Name", Group: "Change Info", Description: "-"},
		{Name: "Change Email", Group: "Change Info", Description: "-"},
		{Name: "Change Birthdate", Group: "Change Info", Description: "-"},
	}

	for _, caseType := range caseTypes {
		if err := db.FirstOrCreate(&caseType, model.CaseTypes{Name: caseType.Name}).Error; err != nil {
			panic("fail to seed case type: " + err.Error())
		}
	}
}
