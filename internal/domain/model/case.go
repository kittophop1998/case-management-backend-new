package model

import (
	"time"

	"github.com/google/uuid"
)

type Cases struct {
	ID                uuid.UUID   `gorm:"primaryKey;default:uuid_generate_v4()" json:"id"`
	Code              string      `gorm:"column:code;type:varchar(100);uniqueIndex" json:"code"`
	CaseTitle         string      `gorm:"column:case_title" json:"caseTitle"`
	CaseTypeID        uuid.UUID   `gorm:"column:case_type_id;type:uuid;default:uuid_generate_v4()" json:"caseTypeId"`
	CaseType          CaseTypes   `gorm:"foreignKey:CaseTypeID;references:ID"`
	QueueID           *uuid.UUID  `gorm:"column:queue_id;type:uuid" json:"queueId"`
	Queue             Queues      `gorm:"foreignKey:QueueID;references:ID"`
	Priority          string      `gorm:"column:priority" json:"priority"`
	CustomerID        string      `gorm:"column:customer_id" json:"customerId"`
	AeonID            string      `gorm:"column:aeon_id" json:"aeonId"`
	CustomerName      string      `gorm:"column:customer_name" json:"customerName"`
	AssignedToUserID  *uuid.UUID  `gorm:"column:assigned_to_user_id;type:uuid" json:"assignedToUserId"`
	AssignedToUser    *User       `gorm:"foreignKey:AssignedToUserID;references:ID" json:"assignedToUser"`
	StatusID          uuid.UUID   `gorm:"column:status_id" json:"statusId"`
	Status            CaseStatus  `gorm:"foreignKey:StatusID;references:ID" json:"status"`
	DispositionMainID *uuid.UUID  `gorm:"column:disposition_main_id;type:uuid" json:"dispositionMainId"`
	DispositionSubID  *uuid.UUID  `gorm:"column:disposition_sub_id;type:uuid" json:"dispositionSubId"`
	StartDate         time.Time   `gorm:"column:start_date;type:date" json:"startDate"`
	DueDate           *time.Time  `gorm:"column:duedate;type:date" json:"dueDate"`
	ClosedDate        time.Time   `gorm:"column:close_date;type:date" json:"closedDate"`
	EndDate           time.Time   `gorm:"column:end_date;type:date" json:"endDate"`
	ProductID         *uuid.UUID  `gorm:"column:product_id;type:uuid" json:"productId"`
	Product           *Products   `gorm:"foreignKey:ProductID;references:ID" json:"product"`
	Channel           *string     `json:"channel"`
	ReasonCodeID      *uuid.UUID  `gorm:"column:reason_code_id;type:uuid" json:"reasonCodeId"`
	ReasonCode        *ReasonCode `gorm:"foreignKey:ReasonCodeID;references:ID" json:"reasonCode"`
	VerifyStatus      *string     `gorm:"column:verify_status" json:"verifyStatus"`
	Description       string      `gorm:"column:description" json:"description"`
	CreatedBy         uuid.UUID   `gorm:"column:created_by;type:uuid" json:"createdBy"`
	Creator           User        `gorm:"foreignKey:CreatedBy;references:ID" json:"creator"`
	CreatedAt         time.Time   `gorm:"column:created_at;type:timestamp" json:"createdAt"`
	UpdatedBy         uuid.UUID   `gorm:"column:updated_by;type:uuid" json:"updatedBy"`
	UpdatedAt         time.Time   `gorm:"column:updated_at;type:timestamp" json:"updatedAt"`
	DeletedBy         uuid.UUID   `gorm:"column:deleted_by;type:uuid" json:"deletedBy"`
	DeletedAt         time.Time   `gorm:"column:deleted_at;type:timestamp" json:"deletedAt"`
}

type CaseDispositionMain struct {
	CaseId            uuid.UUID       `json:"caseId" gorm:"type:uuid"`
	DispositionMainId uuid.UUID       `json:"dispositionMainId" gorm:"type:uuid"`
	Main              DispositionMain `gorm:"foreignKey:DispositionMainId;references:ID"`
	CreatedBy         uuid.UUID       `gorm:"column:created_by;type:uuid" json:"createdBy"`
	CreatedAt         time.Time       `gorm:"column:created_at;type:timestamp" json:"createdAt"`
	UpdatedBy         uuid.UUID       `gorm:"column:updated_by;type:uuid" json:"updatedBy"`
	UpdatedAt         time.Time       `gorm:"column:updated_at;type:timestamp" json:"updatedAt"`
	DeletedBy         uuid.UUID       `gorm:"column:deleted_by;type:uuid" json:"deletedBy"`
	DeletedAt         time.Time       `gorm:"column:deleted_at;type:timestamp" json:"deletedAt"`
}

type CaseDispositionSub struct {
	CaseId           uuid.UUID      `json:"caseId" gorm:"type:uuid"`
	DispositionSubId uuid.UUID      `json:"dispositionSubId" gorm:"type:uuid"`
	Sub              DispositionSub `gorm:"foreignKey:DispositionSubId;references:ID"`
	CreatedBy        uuid.UUID      `gorm:"column:created_by;type:uuid" json:"createdBy"`
	CreatedAt        time.Time      `gorm:"column:created_at;type:timestamp" json:"createdAt"`
	UpdatedBy        uuid.UUID      `gorm:"column:updated_by;type:uuid" json:"updatedBy"`
	UpdatedAt        time.Time      `gorm:"column:updated_at;type:timestamp" json:"updatedAt"`
	DeletedBy        uuid.UUID      `gorm:"column:deleted_by;type:uuid" json:"deletedBy"`
	DeletedAt        time.Time      `gorm:"column:deleted_at;type:timestamp" json:"deletedAt"`
}

type CaseNotes struct {
	ID        uuid.UUID `gorm:"primaryKey;default:uuid_generate_v4()" json:"id"`
	CaseId    uuid.UUID `json:"case_id" gorm:"type:uuid"`
	Content   string    `json:"content" gorm:"type:text"`
	CreatedBy uuid.UUID `json:"created_by" gorm:"type:uuid"`
	Creator   User      `gorm:"foreignKey:CreatedBy;references:ID" json:"creator"`
	CreatedAt time.Time `json:"created_at"`
}

type CaseNotesResponse struct {
	ID        uuid.UUID `json:"id"`
	Content   string    `json:"content"`
	CreatedBy string    `json:"createdBy"`
	CreatedAt time.Time `json:"createdAt"`
}

type VerifyQuestionHistory struct {
	ID               uuid.UUID `json:"id" gorm:"type:uuid;default:uuid_generate_v4()"`
	CustomerId       string    `json:"customer_id"`
	Question         string    `json:"question" gorm:"type:text"`
	AnswerProvided   string    `json:"answer_provided" gorm:"type:text"`
	IsCorrect        bool      `json:"is_correct"`
	VerifyByUserId   uuid.UUID `json:"verify_by_user_id" gorm:"type:uuid;default:uuid_generate_v4()"`
	VerificationDate time.Time `json:"verification_date"`
	CaseId           uuid.UUID `json:"case_id" gorm:"type:uuid;default:uuid_generate_v4()"`
}

type ReasonCode struct {
	ID                uuid.UUID `json:"id" gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	Code              string    `json:"code"`
	DescriptionEn     string    `json:"description_en"`
	DescriptionTh     string    `json:"description_th"`
	Category          string    `json:"category"`
	SLAResponseTime   string    `json:"sla_response_time"`
	SLAResolutionTime string    `json:"sla_resolution_time"`
	Note              string    `json:"note"`
	CreatedBy         uuid.UUID `gorm:"type:uuid" json:"createdBy"`
	CreatedAt         time.Time `gorm:"type:timestamp" json:"createdAt"`
	UpdatedBy         uuid.UUID `gorm:"type:uuid" json:"updatedBy"`
	UpdatedAt         time.Time `gorm:"type:timestamp" json:"updatedAt"`
	DeletedBy         uuid.UUID `gorm:"type:uuid" json:"deletedBy"`
	DeletedAt         time.Time `gorm:"type:timestamp" json:"deletedAt"`
}

// ##### Case Management For Response #####
type CaseResponse struct {
	Code         string `json:"code"`
	CaseID       string `json:"caseId"`
	CustomerID   string `json:"customerId"`
	AeonID       string `json:"aeonId"`
	CustomerName string `json:"customerName"`
	CaseType     string `json:"caseType"`
	CurrentQueue string `json:"currentQueue"`
	CurrentUser  string `json:"currentUser"`
	CaseGroup    string `json:"caseGroup"`
	CasePriority string `json:"casePriority"`
	ClosedDate   string `json:"closedDate"`
	ReceivedFrom string `json:"receivedFrom"`
	SLADate      string `json:"slaDate"`
	Status       string `json:"status"`
	CreatedBy    string `json:"createdBy"`
	CreatedDate  string `json:"createdDate"`
}

type CaseDetailResponse struct {
	Code                string                           `json:"code"`
	CaseType            string                           `json:"caseType"`
	CaseTypeID          string                           `json:"caseTypeId"`
	CaseGroup           string                           `json:"caseGroup"`
	CaseID              string                           `json:"caseId"`
	CreatedBy           string                           `json:"createdBy"`
	CreatedDate         string                           `json:"createdDate"`
	VerifyStatus        *string                          `json:"verifyStatus"`
	Channel             *string                          `json:"channel"`
	Priority            string                           `json:"priority"`
	ReasonCode          *string                          `json:"reasonCode"`
	DueDate             string                           `json:"dueDate"`
	Status              string                           `json:"status"`
	CurrentQueue        string                           `json:"currentQueue"`
	CaseDescription     string                           `json:"caseDescription"`
	AllocateToQueueTeam *string                          `json:"allocateToQueueTeam"`
	Dispositions        []*CaseDispositionDetailResponse `json:"dispositions"`
}

type CaseDispositionDetailResponse struct {
	Main string   `json:"main"`
	Subs []string `json:"subs"`
}

type ReasonCodeResponse struct {
	ID            uuid.UUID `json:"id"`
	Code          string    `json:"code"`
	DescriptionEn string    `json:"descriptionEn"`
	DescriptionTh string    `json:"descriptionTh"`
	Notice        string    `json:"notice"`
}

// ##### Case Management Request #####
type CreateCaseRequest struct {
	CustomerID          string   `json:"customerId"`
	CustomerName        string   `json:"customerName"`
	CaseTypeID          string   `json:"caseTypeId"`
	Channel             *string  `json:"channel"`
	Priority            string   `json:"priority" binding:"required,oneof=High Normal"`
	ReasonCode          *string  `json:"reasonCode"`
	DueDate             *string  `json:"dueDate"`
	DispositionMainID   *string  `json:"dispositionMainId"`
	DispositionMains    []string `json:"dispositionMains"`
	DispositionSubID    *string  `json:"dispositionSubId"`
	DispositionSubs     []string `json:"dispositionSubs"`
	VerifyStatus        *string  `json:"verifyStatus"`
	ProductID           *string  `json:"productId"`
	AllocateToQueueTeam *string  `json:"allocateToQueueTeam"`
	QueueID             *string  `json:"queueId"`
	CaseDescription     string   `json:"caseDescription"`
	CaseNote            []string `json:"caseNote"`
}

type UpdateCaseRequest struct {
	CaseTypeID            *string                `json:"caseTypeId,omitempty"`
	Priority              string                 `json:"priority,omitempty" binding:"required,oneof=High Normal"`
	ReasonCodeID          string                 `json:"reasonCode,omitempty"`
	DueDate               string                 `json:"dueDate,omitempty"`
	ReallocateToQueueTeam string                 `json:"reallocateToQueueTeam,omitempty"`
	Data                  map[string]interface{} `json:"data"`
}

type CaseFilter struct {
	Keyword  string      `form:"keyword" json:"keyword"`
	StatusID []uuid.UUID `form:"statusId" json:"statusId"`
	QueueID  *uuid.UUID  `form:"queueId" json:"queueId"`
	Priority []string    `form:"priority" json:"priority" binding:"omitempty,oneof=High,Normal"`
	Sort     string      `form:"sort" json:"sort"`
}

type CaseQuery struct {
	Category string     `form:"category"`
	Keyword  string     `form:"keyword"`
	Priority []string   `form:"priority"`
	StatusID []string   `form:"statusId"`
	QueueID  *uuid.UUID `form:"queue"`
}

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

type CaseNoteRequest struct {
	Content string `json:"content"`
}

func (Cases) TableName() string {
	return "cases"
}

func (NoteTypes) TableName() string {
	return "note_types"
}

func (CaseNotes) TableName() string {
	return "case_notes"
}

func (ApiLogs) TableName() string {
	return "api_logs"
}

func (CustomerNote) TableName() string {
	return "customers_note"
}

func (c CaseTypes) GetIdentifier() string { return c.Name }
func (c *CaseTypes) GetID() uuid.UUID     { return c.ID }
