package seed

import (
	"case-management/internal/domain/model"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type DispositionMainMap map[string]uuid.UUID

func SeedDispositionMain(db *gorm.DB) DispositionMainMap {
	dispositionMainMap := make(DispositionMainMap)

	dispositionsMain := []model.DispositionMain{
		{Name: "Open Credit Card", Description: "open credit card"},
		{Name: "Close Credit Card", Description: "close credit card"},
		{Name: "Public Information", Description: "public Information"},
	}

	for _, disposition := range dispositionsMain {
		if err := db.FirstOrCreate(&disposition, model.DispositionMain{Name: disposition.Name}).Error; err != nil {
			panic("Failed to seed disposition main, " + err.Error())
		}
		dispositionMainMap[disposition.Name] = disposition.ID
	}
	return dispositionMainMap
}
