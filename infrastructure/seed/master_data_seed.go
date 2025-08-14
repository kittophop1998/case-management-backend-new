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

type RoleMap map[string]uuid.UUID
type PermissionMap map[string]uuid.UUID
type SectionMap map[string]uuid.UUID
type CenterMap map[string]uuid.UUID
type DepartmentMap map[string]uuid.UUID
type QueueMap map[string]uuid.UUID
type DispositionMainMap map[string]uuid.UUID
type CaseTypeMap map[string]uuid.UUID

func SeedMasterData(db *gorm.DB) (RoleMap, SectionMap, CenterMap, DepartmentMap, QueueMap, DispositionMainMap) {
	roleMap := make(RoleMap)
	permissionMap := make(PermissionMap)
	sectionMap := make(SectionMap)
	centerMap := make(CenterMap)
	departmentMap := make(DepartmentMap)
	queueMap := make(QueueMap)
	dispositionMainMap := make(DispositionMainMap)
	caseTypeMap := make(CaseTypeMap)

	seedEntities(db, []Seedable{
		&model.Role{Name: "Admin"},
		&model.Role{Name: "Staff"},
		&model.Role{Name: "Supervisor"},
		&model.Role{Name: "AsstManager Up"},
	}, func(db *gorm.DB, i Seedable) *gorm.DB {
		return db.Where("name = ?", i.GetIdentifier())
	}, func(name string, id uuid.UUID) {
		roleMap[name] = id
	})

	seedEntities(db, []Seedable{
		&model.Permission{Key: "user.login", Name: "Login"},
		&model.Permission{Key: "user.logout", Name: "Logout"},
		&model.Permission{Key: "user.management", Name: "User Management"},
		&model.Permission{Key: "user.profile", Name: "View profile"},
		&model.Permission{Key: "user.assess", Name: "Assess Control"},
		&model.Permission{Key: "user.customersearch", Name: "Customer Search"},
		&model.Permission{Key: "user.verifycustomer", Name: "Verify Customer"},
		&model.Permission{Key: "user.customerdashboard", Name: "Customer Dashboard"},
		&model.Permission{Key: "case.management", Name: "Case Management"},
		&model.Permission{Key: "case.exporthistorical", Name: "Export Case Historical"},
		&model.Permission{Key: "case.standardreport", Name: "Standard Report"},
	}, func(db *gorm.DB, i Seedable) *gorm.DB {
		return db.Where("key = ?", i.GetIdentifier())
	}, func(name string, id uuid.UUID) {
		permissionMap[name] = id
	})

	seedEntities(db, []Seedable{
		&model.Section{Name: "Inbound"},
		&model.Section{Name: "Outbound"},
	}, func(db *gorm.DB, i Seedable) *gorm.DB {
		return db.Where("name = ?", i.GetIdentifier())
	}, func(name string, id uuid.UUID) {
		sectionMap[name] = id
	})

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
		&model.Department{Name: "Marketing"},
		&model.Department{Name: "Sales"},
		&model.Department{Name: "Support"},
	}, func(db *gorm.DB, i Seedable) *gorm.DB {
		return db.Where("name = ?", i.GetIdentifier())
	}, func(name string, id uuid.UUID) {
		departmentMap[name] = id
	})

	seedEntities(db, []Seedable{
		&model.Queue{Name: "Inbound Agent Supervisor BKK Queue"},
		&model.Queue{Name: "Outbound Agent Supervisor HY Queue"},
	}, func(db *gorm.DB, i Seedable) *gorm.DB {
		return db.Where("name = ?", i.GetIdentifier())
	}, func(name string, id uuid.UUID) {
		queueMap[name] = id
	})

	seedEntities(db, []Seedable{
		&model.DispositionMain{Name: "Open Credit Card", Description: "open credit card"},
	}, func(db *gorm.DB, i Seedable) *gorm.DB {
		return db.Where("name = ?", i.GetIdentifier())
	}, func(name string, id uuid.UUID) {
		dispositionMainMap[name] = id
	})

	seedEntities(db, []Seedable{
		&model.CaseTypes{Name: "Inquiry And Disposition"},
	}, func(db *gorm.DB, i Seedable) *gorm.DB {
		return db.Where("name = ?", i.GetIdentifier())
	}, func(name string, id uuid.UUID) {
		caseTypeMap[name] = id
	})

	return roleMap, sectionMap, centerMap, departmentMap, queueMap, dispositionMainMap
}
