package seed

import (
	"case-management/internal/domain/model"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type DispositionSubMap map[string]uuid.UUID

func SeedDispositionSub(db *gorm.DB, dispositionMain map[string]uuid.UUID) DispositionSubMap {
	dispositionSubMap := make(DispositionSubMap)

	dispositionSub := []model.DispositionSub{
		{Name: "Aeon", MainID: dispositionMain["Open Credit Card"]},
		{Name: "Aeon Thai Smile", MainID: dispositionMain["Open Credit Card"]},
	}

	for _, dispoSub := range dispositionSub {
		if err := db.FirstOrCreate(&dispoSub, model.DispositionSub{Name: dispoSub.Name}).Error; err != nil {
			panic("Failed to seed disposition sub: " + err.Error())
		}
		dispositionSubMap[dispoSub.Name] = dispoSub.ID
	}

	return dispositionSubMap
}
