package model

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID           uuid.UUID  `json:"id" gorm:"type:uuid;default:uuid_generate_v4()"`
	StaffID      *uint      `json:"staffId" gorm:"uniqueIndex"`
	Username     string     `gorm:"uniqueIndex;type:varchar(50)" json:"username" example:"user"`
	Email        string     `gorm:"type:varchar(100)" json:"email" example:"user@example.com"`
	Name         string     `gorm:"type:varchar(100)" json:"name"`
	SectionID    uuid.UUID  `json:"sectionId"`
	Section      Section    `gorm:"foreignKey:SectionID" json:"section"`
	IsActive     *bool      `json:"isActive" gorm:"default:true"`
	CenterID     uuid.UUID  `json:"centerId"`
	Center       Center     `json:"center"`
	RoleID       uuid.UUID  `json:"roleId"`
	Role         Role       `json:"role"`
	OperatorID   *uint      `json:"operatorId"`
	Department   Department `json:"department"`
	DepartmentID uuid.UUID  `json:"departmentId"`
	Password     string     `json:"-"`
	UserTypes    string     `json:"-"`
	CreatedAt    time.Time  `json:"createdAt"`
	CreatedBy    uuid.UUID  `json:"createdBy"`
	UpdatedAt    time.Time  `json:"updatedAt"`
	UpdatedBy    uuid.UUID  `json:"updatedBy"`
}
