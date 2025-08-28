package seed

import (
	"case-management/internal/domain/model"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type CenterMap map[string]uuid.UUID

func SeedCenter(db *gorm.DB) CenterMap {
	centerMap := make(CenterMap)

	centers := []model.Center{
		{Name: "BKK"},
		{Name: "HY"},
		{Name: "KK"},
	}

	for _, center := range centers {
		if err := db.FirstOrCreate(&center, model.Center{Name: center.Name}).Error; err != nil {
			panic("Failed to seed center: " + err.Error())
		}
		centerMap[center.Name] = center.ID
	}
	return centerMap
}
