package model

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID            uuid.UUID  `gorm:"column:id;primaryKey;type:uuid;default:uuid_generate_v4()" json:"id" `
	StaffID       *uint      `gorm:"column:staff_id;uniqueIndex" json:"staffId"`
	Username      string     `gorm:"column:username;uniqueIndex" json:"username" example:"user"`
	Email         string     `gorm:"column:email;type:varchar(100)" json:"email" example:"user@example.com"`
	Name          string     `gorm:"column:name;type:varchar(100)" json:"name"`
	SectionID     uuid.UUID  `gorm:"column:section_id" json:"sectionId"`
	Section       Section    `gorm:"foreignKey:SectionID" json:"section"`
	IsActive      *bool      `gorm:"column:is_active;default:true" json:"isActive"`
	CenterID      uuid.UUID  `gorm:"column:center_id" json:"centerId"`
	Center        Center     `gorm:"foreignKey:CenterID" json:"center"`
	RoleID        uuid.UUID  `gorm:"column:role_id" json:"roleId"`
	Role          Role       `json:"role"`
	DepartmentID  uuid.UUID  `gorm:"column:department_id" json:"departmentId"`
	Department    Department `json:"department"`
	OperatorID    *uint      `gorm:"column:operator_id" json:"operatorId"`
	Password      string     `json:"-"`
	UserTypes     string     `json:"-"`
	CreatedAt     time.Time  `gorm:"column:created_at" json:"createdAt"`
	CreatedUserBy uuid.UUID  `gorm:"column:created_by;type:uuid" json:"createdBy"`
	UpdatedAt     time.Time  `gorm:"column:updated_at" json:"updatedAt"`
	UpdatedBy     uuid.UUID  `gorm:"column:updated_by;type:uuid" json:"updatedBy"`
}

// ##### User For Request #####
type UserFilter struct {
	Keyword      string    `json:"keyword"`
	Name         string    `gorm:"type:varchar(100)" json:"name"`
	Sort         string    `json:"sort"`
	IsActive     *bool     `json:"isActive"`
	Role         string    `json:"role"`
	Section      Section   `json:"section"`
	Center       string    `json:"center"`
	RoleID       uuid.UUID `json:"roleID,omitempty"`
	SectionID    uuid.UUID `json:"sectionID,omitempty"`
	CenterID     uuid.UUID `json:"centerID,omitempty"`
	DepartmentID uuid.UUID `json:"departmentID,omitempty"`
	QueueID      uuid.UUID `json:"queueID,omitempty"`
	IsNotInQueue *bool     `json:"isNotInQueue,omitempty"`
}

// ##### User For Response #####
type UserProfileResponse struct {
	ID          uuid.UUID    `json:"id"`
	Username    string       `json:"username"`
	Name        string       `json:"name"`
	Email       string       `json:"email"`
	Role        Role         `json:"role"`
	Department  Department   `json:"department"`
	Section     Section      `json:"section"`
	Center      Center       `json:"center"`
	Permissions []Permission `json:"permissions"`
}
