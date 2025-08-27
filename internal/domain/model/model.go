package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

// ##### Case Management #####
type CaseFilter struct {
	Keyword     string     `form:"keyword" json:"keyword"`
	StatusID    *uint      `form:"statusId" json:"statusId"`
	PriorityID  *uint      `form:"priorityId" json:"priorityId"`
	SLADateFrom *time.Time `form:"slaDateFrom" json:"slaDateFrom"`
	SLADateTo   *time.Time `form:"slaDateTo" json:"slaDateTo"`
	Sort        string     `form:"sort" json:"sort"`
}

type CreateCaseRequest struct {
	CustomerId        string         `json:"customerId" binding:"required"`
	CaseTypeId        uuid.UUID      `json:"caseTypeId" binding:"required"`
	DispositionMainId uuid.UUID      `json:"dispositionMainId" gorm:"type:uuid" binding:"required"`
	DispositionMains  datatypes.JSON `json:"dispositionMains" gorm:"type:jsonb" binding:"required"`
	CaseDescription   string         `json:"caseDescription" gorm:"type:text"`
	CaseNote          datatypes.JSON `json:"caseNote" gorm:"type:jsonb"`
}

// ##### User Management #####
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

// ##### Customer #####

type StatusResponse struct {
	Status string `json:"status"`
}

type MessageResponse struct {
	Message string `json:"message"`
}

type AccessTokenRequest struct {
	Access_token string `json:"access_token" binding:"required"`
}

type PermissionWithRolesResponse struct {
	Permission string   `json:"permission"`
	Name       string   `json:"name"`
	Roles      []string `json:"roles"`
}

type AccessControls struct {
	Permissions []PermissionWithRolesResponse `json:"permissions"`
	PermCount   int                           `json:"permCount"`
}

type FormFilter struct {
	Limit  int                    `json:"limit"`
	Page   int                    `json:"page"`
	Sort   string                 `json:"sort"`
	Filter map[string]interface{} `json:"filter"`
}

type CreateNoteRequest struct {
	NoteTypeID string `json:"noteTypeId" binding:"required"`
	Note       string `json:"note" binding:"required"`
}

type Model struct {
	ID        uuid.UUID      `gorm:"primaryKey;default:uuid_generate_v4()" json:"id"`
	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-" swaggerignore:"true"`
}
