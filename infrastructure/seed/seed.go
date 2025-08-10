package seed

import (
	"log"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

func SeedAllData(db *gorm.DB) error {
	// Seed master data
	roleMap, sectionMap, centerMap, departmentMap, queueMap, dispositionMainMap := SeedMasterData(db)

	// Seed disposition sub data
	SeedDispositionSub(db, dispositionMainMap)

	// Seed role permissions
	// if err := SeedRolePermission(db); err != nil {
	// 	return err
	// }

	// Seed users
	SeedUser(db, roleMap, sectionMap, centerMap, departmentMap, queueMap)

	return nil
}

func seedEntities(
	db *gorm.DB,
	items []Seedable,
	whereBuilder func(*gorm.DB, Seedable) *gorm.DB,
	mapSetter func(string, uuid.UUID),
) {
	for _, item := range items {
		key := item.GetIdentifier()
		if key == "" {
			log.Println("Skipping empty identifier")
			continue
		}

		query := whereBuilder(db, item)
		if err := query.FirstOrCreate(item).Error; err != nil {
			log.Printf("Failed to seed %T %s: %v", item, key, err)
		} else {
			mapSetter(key, item.GetID())
		}
	}
}
