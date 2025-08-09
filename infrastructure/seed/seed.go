package seed

import "gorm.io/gorm"

func SeedAllData(db *gorm.DB) error {
	// Seed master data
	roleMap, sectionMap, centerMap, departmentMap, queueMap := SeedMasterData(db)

	// Seed role permissions
	if err := SeedRolePermission(db); err != nil {
		return err
	}

	// Seed users
	SeedUser(db, roleMap, sectionMap, centerMap, departmentMap, queueMap)

	return nil
}
