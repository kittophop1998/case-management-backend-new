package seed

import (
	"case-management/internal/domain/model"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type SectionMap map[string]uuid.UUID

func SeedSection(db *gorm.DB) SectionMap {
	sectionMap := make(SectionMap)

	sections := []model.Section{
		{Key: "SYSTEM", Name: "System"},
		{Key: "CHB", Name: "Chargeback"},
		{Key: "CSD", Name: "Customer service Development"},
		{Key: "ONB", Name: "Inbound BBK"},
		{Key: "ONC", Name: "Inbound CM"},
		{Key: "ONK", Name: "Inbound KK"},
		{Key: "ONH", Name: "Inbound HY"},
	}

	for _, section := range sections {
		if err := db.FirstOrCreate(&section, model.Section{Name: section.Name}).Error; err != nil {
			panic(err)
		}
		sectionMap[section.Key] = section.ID
	}

	return sectionMap
}
