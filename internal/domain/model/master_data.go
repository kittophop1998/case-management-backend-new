package model

import "github.com/google/uuid"

type Permission struct {
	ID    uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey" json:"id"`
	Key   string    `gorm:"uniqueIndex;not null" json:"key"`
	Name  string    `gorm:"uniqueIndex;not null" json:"name"`
	Roles []Role    `gorm:"many2many:role_permissions"`
}

type Role struct {
	ID          uuid.UUID    `gorm:"type:uuid;default:uuid_generate_v4();primaryKey" json:"id"`
	Name        string       `gorm:"uniqueIndex;not null" json:"name"`
	Permissions []Permission `gorm:"many2many:role_permissions" json:"permissions"`
}

type RolePermission struct {
	RoleID       uuid.UUID `gorm:"type:uuid;not null;primaryKey"`
	PermissionID uuid.UUID `gorm:"type:uuid;not null;primaryKey"`
}

type Center struct {
	ID   uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey" json:"id"`
	Name string    `gorm:"type:varchar(100)" json:"name"`
}

type Team struct {
	ID   uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey" json:"id"`
	Name string    `gorm:"type:varchar(100)" json:"name"`
}

type Queue struct {
	ID   uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey" json:"id"`
	Name string    `gorm:"type:varchar(100)" json:"name"`
}

type Department struct {
	ID   uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey" json:"id"`
	Name string    `gorm:"type:varchar(100)" json:"name"`
}

type AddInitialDescriptionRequest struct {
	CaseID      string `json:"case_id" binding:"required,uuid"`
	Description string `json:"description" binding:"required"`
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

func (Team) TableName() string {
	return "teams"
}

func (Department) TableName() string {
	return "departments"
}
