package seed

import (
	"gorm.io/gorm"
)

func SeedAllData(db *gorm.DB) error {
	centerMap := SeedCenter(db)
	roleMap := SeedRole(db)
	permissionMap := SeedPermission(db)
	departmentMap := SeedDepartment(db)
	sectionMap := SeedSection(db)
	SeedNoteTypes(db)
	SeedCaseType(db)
	SeedRolePermission(db, roleMap, permissionMap, departmentMap, sectionMap)
	SeedUser(db, roleMap, sectionMap, centerMap, departmentMap)
	dispositionMainMap := SeedDispositionMain(db)
	SeedDispositionSub(db, dispositionMainMap)
	SeedCaseStatus(db)

	return nil
}
