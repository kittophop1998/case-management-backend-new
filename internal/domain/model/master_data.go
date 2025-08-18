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

type Queue struct {
	ID   uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey" json:"id"`
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

type DispositionMain struct {
	ID          uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey" json:"id"`
	Name        string    `gorm:"type:varchar(100)" json:"name"`
	Description string    `gorm:"type:text" json:"description"`
	// Subs        []DispositionSub `json:"subs" gorm:"foreignKey:MainID;references:ID"`
}

type DispositionSub struct {
	ID          uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey" json:"id"`
	MainID      uuid.UUID `gorm:"type:uuid;not null" json:"main_id"`
	Name        string    `gorm:"type:varchar(100)" json:"name"`
	Description string    `gorm:"type:text" json:"description"`
}

type DispositionFilter struct {
	Keyword string `form:"keyword" json:"keyword"`
	Limit   int    `form:"limit" json:"limit"`
	Offset  int    `form:"offset" json:"offset"`
}

func (r Role) GetIdentifier() string { return r.Name }
func (r *Role) GetID() uuid.UUID     { return r.ID }

func (p Permission) GetIdentifier() string { return p.Key }
func (p *Permission) GetID() uuid.UUID     { return p.ID }

func (s Section) GetIdentifier() string { return s.Name }
func (s *Section) GetID() uuid.UUID     { return s.ID }

func (c Center) GetIdentifier() string { return c.Name }
func (c *Center) GetID() uuid.UUID     { return c.ID }

func (d Department) GetIdentifier() string { return d.Name }
func (d *Department) GetID() uuid.UUID     { return d.ID }

func (q Queue) GetIdentifier() string { return q.Name }
func (q *Queue) GetID() uuid.UUID     { return q.ID }

func (d DispositionMain) GetIdentifier() string { return d.Name }
func (d *DispositionMain) GetID() uuid.UUID     { return d.ID }

func (d DispositionSub) GetIdentifier() string { return d.Name }
func (d *DispositionSub) GetID() uuid.UUID     { return d.ID }

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

func (Queue) TableName() string {
	return "queues"
}

func (DispositionMain) TableName() string {
	return "disposition_mains"
}

func (DispositionSub) TableName() string {
	return "disposition_subs"
}
