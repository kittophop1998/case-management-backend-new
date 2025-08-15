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
		{Name: "System"},
		{Name: "CHL"},
	}

	for _, section := range sections {
		if err := db.FirstOrCreate(&section, model.Section{Name: section.Name}).Error; err != nil {
			panic(err)
		}
		sectionMap[section.Name] = section.ID
	}

	return sectionMap
}
