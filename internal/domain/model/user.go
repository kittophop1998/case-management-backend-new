package model

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID           uuid.UUID  `gorm:"column:id;primaryKey;type:uuid;default:uuid_generate_v4()" json:"id" `
	StaffID      *uint      `gorm:"column:staff_id;uniqueIndex" json:"staffId"`
	Username     string     `gorm:"column:username;uniqueIndex" json:"username" example:"user"`
	Email        string     `gorm:"column:email;type:varchar(100)" json:"email" example:"user@example.com"`
	Name         string     `gorm:"column:name;type:varchar(100)" json:"name"`
	SectionID    uuid.UUID  `gorm:"column:section_id" json:"sectionId"`
	Section      Section    `gorm:"foreignKey:SectionID" json:"section"`
	IsActive     *bool      `gorm:"column:is_active;default:true" json:"isActive"`
	CenterID     uuid.UUID  `gorm:"column:center_id" json:"centerId"`
	Center       Center     `gorm:"foreignKey:CenterID" json:"center"`
	RoleID       uuid.UUID  `gorm:"column:role_id" json:"roleId"`
	Role         Role       `json:"role"`
	DepartmentID uuid.UUID  `gorm:"column:department_id" json:"departmentId"`
	Department   Department `json:"department"`
	OperatorID   *uint      `gorm:"column:operator_id" json:"operatorId"`
	Password     string     `json:"-"`
	UserTypes    string     `json:"-"`
	CreatedAt    time.Time  `gorm:"column:created_at" json:"createdAt"`
	CreatedBy    uuid.UUID  `gorm:"column:created_by_user" json:"createdBy"`
	UpdatedAt    time.Time  `gorm:"column:updated_at" json:"updatedAt"`
	UpdatedBy    uuid.UUID  `gorm:"column:updated_by_user" json:"updatedBy"`
}
