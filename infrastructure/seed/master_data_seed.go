package seed

import (
	"case-management/internal/domain/model"
	"log"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Common interface for seedable entities
type Seedable interface {
	GetIdentifier() string // Name or Key
	GetID() uuid.UUID
}

type RoleMap map[string]uuid.UUID
type PermissionMap map[string]uuid.UUID
type SectionMap map[string]uuid.UUID
type CenterMap map[string]uuid.UUID
type DepartmentMap map[string]uuid.UUID
type QueueMap map[string]uuid.UUID

func SeedMasterData(db *gorm.DB) (RoleMap, SectionMap, CenterMap, DepartmentMap, QueueMap) {
	roleMap := make(RoleMap)
	sectionMap := make(SectionMap)
	centerMap := make(CenterMap)
	departmentMap := make(DepartmentMap)
	queueMap := make(QueueMap)

	seedEntities(db, []Seedable{
		&model.Role{Name: "Admin"},
		&model.Role{Name: "Staff"},
		&model.Role{Name: "Supervisor"},
		&model.Role{Name: "AsstManager Up"},
	}, func(name string, id uuid.UUID) { roleMap[name] = id })

	seedEntities(db, []Seedable{
		&model.Section{Name: "Inbound"},
		&model.Section{Name: "Outbound"},
	}, func(name string, id uuid.UUID) { sectionMap[name] = id })

	seedEntities(db, []Seedable{
		&model.Center{Name: "BKK"},
		&model.Center{Name: "HY"},
		&model.Center{Name: "KK"},
	}, func(name string, id uuid.UUID) { centerMap[name] = id })

	seedEntities(db, []Seedable{
		&model.Department{Name: "Marketing"},
		&model.Department{Name: "Sales"},
		&model.Department{Name: "Support"},
	}, func(name string, id uuid.UUID) { departmentMap[name] = id })

	seedEntities(db, []Seedable{
		&model.Queue{Name: "Inbound Agent Supervisor BKK Queue"},
		&model.Queue{Name: "Outbound Agent Supervisor HY Queue"},
	}, func(name string, id uuid.UUID) { queueMap[name] = id })

	return roleMap, sectionMap, centerMap, departmentMap, queueMap
}

func seedEntities(db *gorm.DB, items []Seedable, mapSetter func(string, uuid.UUID)) {
	for _, item := range items {
		key := item.GetIdentifier()
		if key == "" {
			log.Println("Skipping empty identifier")
			continue
		}
		if err := db.Where("name = ?", key).FirstOrCreate(item).Error; err != nil {
			log.Printf("Failed to seed %T %s: %v", item, key, err)
		} else {
			mapSetter(key, item.GetID())
		}
	}
}
