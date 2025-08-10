package seed

import (
	"case-management/internal/domain/model"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type DispositionSubMap map[string]uuid.UUID

func SeedDispositionSub(db *gorm.DB, dispositionMain map[string]uuid.UUID) DispositionSubMap {
	dispositionSubMap := make(DispositionSubMap)

	seedEntities(db, []Seedable{
		&model.DispositionSub{Name: "Loan Inquiry", MainID: dispositionMain["Check loan"]},
		&model.DispositionSub{Name: "Product Inquiry", MainID: dispositionMain["Inquire about product"]},
	}, func(db *gorm.DB, i Seedable) *gorm.DB {
		return db.Where("name = ?", i.GetIdentifier())
	}, func(name string, id uuid.UUID) {
		dispositionSubMap[name] = id
	})

	return dispositionSubMap
}
