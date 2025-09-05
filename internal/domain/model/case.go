package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"
)

type Cases struct {
	ID                uuid.UUID  `gorm:"primaryKey;default:uuid_generate_v4()" json:"id"`
	CaseTitle         string     `json:"caseTitle"`
	CaseTypeID        uuid.UUID  `gorm:"type:uuid;default:uuid_generate_v4()" json:"caseTypeId"`
	CaseType          CaseTypes  `gorm:"foreignKey:CaseTypeID;references:ID"`
	QueueID           *uuid.UUID `json:"queueId" gorm:"type:uuid"`
	Queue             Queues     `json:"queue"`
	Priority          string     `json:"priority"`
	CustomerID        string     `json:"customerId"`
	AeonID            string     `json:"aeonId"`
	CustomerName      string     `json:"customerName"`
	AssignedToUserID  *uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4()" json:"assignedToUserId"`
	AssignedToUser    *User      `gorm:"foreignKey:AssignedToUserID;references:ID" json:"assignedToUser"`
	StatusID          uuid.UUID  `json:"statusId"`
	Status            CaseStatus `gorm:"foreignKey:StatusID;references:ID"`
	DispositionMainID *uuid.UUID `gorm:"type:uuid" json:"dispositionMainId"`
	DispositionSubID  *uuid.UUID `gorm:"type:uuid" json:"dispositionSubId"`
	StartDate         time.Time  `gorm:"type:date" json:"startDate"`
	DueDate           *time.Time `gorm:"type:date" json:"dueDate"`
	ClosedDate        time.Time  `gorm:"type:date" json:"closedDate"`
	EndDate           time.Time  `gorm:"type:date" json:"endDate"`
	ProductID         *uuid.UUID `gorm:"type:uuid" json:"productId"`
	Product           *Products  `gorm:"foreignKey:ProductID;references:ID" json:"product"`
	Channel           *string    `json:"channel"`
	ReasonCode        *string    `json:"reasonCode"`
	VerifyStatus      *string    `json:"verifyStatus"`
	Description       string     `json:"description"`
	CreatedBy         uuid.UUID  `gorm:"type:uuid" json:"createdBy"`
	Creator           User       `gorm:"foreignKey:CreatedBy;references:ID" json:"creator"`
	CreatedAt         time.Time  `gorm:"type:timestamp" json:"createdAt"`
	UpdatedBy         uuid.UUID  `gorm:"type:uuid" json:"updatedBy"`
	UpdatedAt         time.Time  `gorm:"type:timestamp" json:"updatedAt"`
	DeletedBy         uuid.UUID  `gorm:"type:uuid" json:"deletedBy"`
	DeletedAt         time.Time  `gorm:"type:timestamp" json:"deletedAt"`
}

type CaseDispositionMain struct {
	Model
	CaseId            uuid.UUID `json:"caseId" gorm:"type:uuid"`
	DispositionMainId uuid.UUID `json:"dispositionMainId" gorm:"type:uuid"`
}

type CaseDispositionSub struct {
	Model
	CaseId           uuid.UUID `json:"caseId" gorm:"type:uuid"`
	DispositionSubId uuid.UUID `json:"dispositionSubId" gorm:"type:uuid"`
}

type CaseNotes struct {
	ID        uuid.UUID `gorm:"primaryKey;default:uuid_generate_v4()" json:"id"`
	CaseId    uuid.UUID `json:"case_id" gorm:"type:uuid"`
	UserId    uuid.UUID `json:"user_id" gorm:"type:uuid"`
	Content   string    `json:"content" gorm:"type:text"`
	CreatedAt time.Time `json:"created_at"`
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
	ID                uuid.UUID `json:"id" gorm:"type:uuid;default:uuid_generate_v4()"`
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
	CaseID       string `json:"caseId"`
	CustomerID   string `json:"customerId"`
	AeonID       string `json:"aeonId"`
	CustomerName string `json:"customerName"`
	CaseType     string `json:"caseType"`
	CurrentQueue string `json:"currentQueue"`
	CurrentUser  string `json:"currentUser"`
	CreateDate   string `json:"createDate"`
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
	CaseType            string  `json:"caseType"`
	CaseID              string  `json:"caseId"`
	CreatedBy           string  `json:"createdBy"`
	CreatedDate         string  `json:"createdDate"`
	VerifyStatus        *string `json:"verifyStatus"`
	Channel             *string `json:"channel"`
	Priority            string  `json:"priority"`
	ReasonCode          *string `json:"reasonCode"`
	DueDate             string  `json:"dueDate"`
	Status              string  `json:"status"`
	CurrentQueue        string  `json:"currentQueue"`
	CaseDescription     string  `json:"caseDescription"`
	AllocateToQueueTeam *string `json:"allocateToQueueTeam"`
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
	CustomerID          string         `json:"customerId"`
	CustomerName        string         `json:"customerName"`
	CaseTypeID          string         `json:"caseTypeId"`
	Channel             *string        `json:"channel"`
	Priority            string         `json:"priority"`
	ReasonCode          *string        `json:"reasonCode"`
	DueDate             *string        `json:"dueDate"`
	DispositionMainID   *string        `json:"dispositionMainId"`
	DispositionMains    datatypes.JSON `json:"dispositionMains"`
	DispositionSubID    *string        `json:"dispositionSubId"`
	DispositionSubs     datatypes.JSON `json:"dispositionSubs"`
	VerifyStatus        *string        `json:"verifyStatus"`
	ProductID           *string        `json:"productId"`
	AllocateToQueueTeam *string        `json:"allocateToQueueTeam"`
	QueueID             *string        `json:"queueId"`
	CaseDescription     string         `json:"caseDescription"`
	CaseNote            datatypes.JSON `json:"caseNote"`
}

type CaseFilter struct {
	Keyword     string     `form:"keyword" json:"keyword"`
	StatusID    *uint      `form:"statusId" json:"statusId"`
	PriorityID  *uint      `form:"priorityId" json:"priorityId"`
	SLADateFrom *time.Time `form:"slaDateFrom" json:"slaDateFrom"`
	SLADateTo   *time.Time `form:"slaDateTo" json:"slaDateTo"`
	Sort        string     `form:"sort" json:"sort"`
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
