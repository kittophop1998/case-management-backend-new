package model

import "github.com/google/uuid"

type User struct {
	Model
	StaffID      *uint      `json:"staffId" gorm:"unique"`
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
	Password     string     `json:"password"`
	UserTypes    string     `json:"userTypes"`
}
