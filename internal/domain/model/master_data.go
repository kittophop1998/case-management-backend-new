package model

import "github.com/google/uuid"

type Permission struct {
	ID   uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey" json:"id"`
	Key  string    `gorm:"uniqueIndex;not null" json:"key"`
	Name string    `gorm:"uniqueIndex;not null" json:"name"`
}

type Role struct {
	ID   uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey" json:"id"`
	Name string    `gorm:"uniqueIndex;not null" json:"name"`
}

// join table
type RolePermission struct {
	RoleID       uuid.UUID `gorm:"type:uuid;not null;primaryKey"`
	PermissionID uuid.UUID `gorm:"type:uuid;not null;primaryKey"`
	SectionID    uuid.UUID `gorm:"type:uuid;not null;primaryKey"`
	DepartmentID uuid.UUID `gorm:"type:uuid;not null;primaryKey"`
}

type Center struct {
	ID   uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey" json:"id"`
	Name string    `gorm:"type:varchar(100)" json:"name"`
}

type Section struct {
	ID   uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey" json:"id"`
	Key  string    `gorm:"uniqueIndex;not null" json:"key"`
	Name string    `gorm:"type:varchar(100)" json:"name"`
}

type Department struct {
	ID   uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey" json:"id"`
	Key  string    `gorm:"uniqueIndex;not null" json:"key"`
	Name string    `gorm:"type:varchar(100)" json:"name"`
}

type AddInitialDescriptionRequest struct {
	CaseID      string `json:"case_id" binding:"required,uuid"`
	Description string `json:"description" binding:"required"`
}

type Products struct {
	ID   uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey" json:"id"`
	Name string    `gorm:"type:varchar(100);not null" json:"name"`
}

type DispositionFilter struct {
	Keyword string `form:"keyword" json:"keyword"`
	Limit   int    `form:"limit" json:"limit"`
	Offset  int    `form:"offset" json:"offset"`
}

type DispositionMain struct {
	ID          uuid.UUID        `gorm:"type:uuid;default:uuid_generate_v4();primaryKey" json:"id"`
	NameTH      string           `gorm:"type:varchar(100)" json:"nameTh"`
	NameEN      string           `gorm:"type:varchar(100)" json:"nameEn"`
	Description string           `gorm:"type:text" json:"description"`
	Subs        []DispositionSub `json:"subs" gorm:"foreignKey:MainID;references:ID"`
}

type DispositionSub struct {
	ID          uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey" json:"id"`
	MainID      uuid.UUID `gorm:"type:uuid;not null" json:"main_id"`
	NameTH      string    `gorm:"type:varchar(100)" json:"nameTh"`
	NameEN      string    `gorm:"type:varchar(100)" json:"nameEn"`
	Description string    `gorm:"type:text" json:"description"`
}

// Disposition Response
type DispositionItem struct {
	DispositionMain DispositionMainRes  `json:"dispositionMain"`
	DispositionSubs []DispositionSubRes `json:"dispositionSubs"`
}

type DispositionMainRes struct {
	ID string `json:"id"`
	TH string `json:"th"`
	EN string `json:"en"`
}

type DispositionSubRes struct {
	ID   string     `json:"id"`
	Name SubNameRes `json:"name"`
}

type SubNameRes struct {
	TH string `json:"th"`
	EN string `json:"en"`
}

func (RolePermission) TableName() string {
	return "role_permissions"
}

func (Permission) TableName() string {
	return "permissions"
}

func (Role) TableName() string {
	return "roles"
}

func (Center) TableName() string {
	return "centers"
}

func (Section) TableName() string {
	return "sections"
}

func (Department) TableName() string {
	return "departments"
}

func (DispositionMain) TableName() string {
	return "disposition_mains"
}

func (DispositionSub) TableName() string {
	return "disposition_subs"
}
