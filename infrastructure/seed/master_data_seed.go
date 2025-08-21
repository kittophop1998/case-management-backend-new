package seed

import (
	"case-management/internal/domain/model"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Common interface for seedable entities
type Seedable interface {
	GetIdentifier() string // Name or Key
	GetID() uuid.UUID
}

type CenterMap map[string]uuid.UUID
type QueueMap map[string]uuid.UUID

func SeedMasterData(db *gorm.DB) (CenterMap, QueueMap) {
	centerMap := make(CenterMap)
	queueMap := make(QueueMap)

	seedEntities(db, []Seedable{
		&model.Center{Name: "BKK"},
		&model.Center{Name: "HY"},
		&model.Center{Name: "KK"},
	}, func(db *gorm.DB, i Seedable) *gorm.DB {
		return db.Where("name = ?", i.GetIdentifier())
	}, func(name string, id uuid.UUID) {
		centerMap[name] = id
	})

	seedEntities(db, []Seedable{
		&model.Queue{Name: "Inbound Agent Supervisor BKK Queue"},
		&model.Queue{Name: "Outbound Agent Supervisor HY Queue"},
	}, func(db *gorm.DB, i Seedable) *gorm.DB {
		return db.Where("name = ?", i.GetIdentifier())
	}, func(name string, id uuid.UUID) {
		queueMap[name] = id
	})

	return centerMap, queueMap
}
