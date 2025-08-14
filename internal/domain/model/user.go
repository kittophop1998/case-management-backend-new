package model

import "github.com/google/uuid"

type User struct {
	Model
	StaffID      *uint      `json:"staffId"`
	Username     string     `gorm:"type:varchar(50)" json:"username" example:"user"`
	Email        string     `gorm:"type:varchar(100)" json:"email" example:"user@example.com"`
	Name         string     `gorm:"type:varchar(100)" json:"name"`
	SectionID    uuid.UUID  `json:"sectionId"`
	Section      Section    `gorm:"foreignKey:SectionID" json:"section"`
	IsActive     *bool      `json:"isActive" gorm:"default:true"`
	CenterID     uuid.UUID  `json:"centerId"`
	Center       Center     `gorm:"foreignKey:CenterID" json:"center"`
	RoleID       uuid.UUID  `json:"roleId"`
	Role         Role       `gorm:"foreignKey:RoleID" json:"role"`
	OperatorID   *uint      `json:"operatorId"`
	Department   Department `gorm:"foreignKey:DepartmentID" json:"department"`
	DepartmentID uuid.UUID  `json:"departmentId"`
}

// User UseCase
type CreateUpdateUserRequest struct {
	StaffID      *uint     `gorm:"uniqueIndex" json:"staffId" validate:"required" example:"12337"`
	Name         string    `json:"name" validate:"required" example:"Janet Adebayo"`
	Username     string    `gorm:"type:varchar(50)" json:"username" example:"user"`
	Email        string    `json:"email" validate:"required" example:"Janet@exam.com"`
	SectionID    uuid.UUID `json:"sectionId" validate:"required" example:"b94eee08-8324-4d4f-b166-d82775553a7e"`
	OperatorID   *uint     `json:"operatorId" validate:"required" example:"1233"`
	CenterID     uuid.UUID `json:"centerId" validate:"required" example:"b94eee08-8324-4d4f-b166-d82775553a7e"`
	RoleID       uuid.UUID `json:"roleId" validate:"required" example:"538cd6c5-4cb3-4463-b7d5-ac6645815476"`
	DepartmentID uuid.UUID `json:"departmentId" validate:"required" example:"b94eee08-8324-4d4f-b166-d82775553a7e"`
	IsActive     *bool     `json:"isActive" validate:"required" example:"true"`
}

type UserFilter struct {
	Keyword   string    `json:"keyword"`
	Name      string    `gorm:"type:varchar(100)" json:"name"`
	Sort      string    `json:"sort"`
	IsActive  *bool     `json:"isActive"`
	Role      string    `json:"role"`
	Section   Section   `json:"section"`
	Center    string    `json:"center"`
	RoleID    uuid.UUID `json:"roleID,omitempty" gorm:"type:uuid;default:uuid_generate_v4()"`
	SectionID uuid.UUID `json:"sectionID,omitempty" gorm:"type:uuid;default:uuid_generate_v4()"`
	CenterID  uuid.UUID `json:"centerID,omitempty" gorm:"type:uuid;default:uuid_generate_v4()"`
}
