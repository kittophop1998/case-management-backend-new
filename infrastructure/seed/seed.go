package seed

import (
	"log"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

func SeedAllData(db *gorm.DB) error {
	centerMap, queueMap, dispositionMainMap := SeedMasterData(db)
	SeedNoteTypes(db)
	SeedCaseType(db)
	roleMap := SeedRole(db)
	permissionMap := SeedPermission(db)
	departmentMap := SeedDepartment(db)
	sectionMap := SeedSection(db)
	SeedRolePermission(db, roleMap, permissionMap, departmentMap, sectionMap)
	SeedUser(db, roleMap, sectionMap, centerMap, departmentMap, queueMap)
	SeedDispositionSub(db, dispositionMainMap)

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
