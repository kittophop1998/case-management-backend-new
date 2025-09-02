package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"
)

type Cases struct {
	ID                uuid.UUID  `gorm:"primaryKey;default:uuid_generate_v4()" json:"id"`
	CaseTitle         string     `json:"caseTitle"`
	CaseTypeID        uuid.UUID  `json:"caseTypeId" gorm:"type:uuid;default:uuid_generate_v4()"`
	CaseType          CaseTypes  `gorm:"foreignKey:CaseTypeID;references:ID"`
	QueueID           *uuid.UUID `json:"queueId" gorm:"type:uuid"`
	Queue             Queues     `json:"queue"`
	Priority          string     `json:"priority"`
	CustomerID        string     `json:"customerId"`
	AeonID            string     `json:"aeonId"`
	CustomerName      string     `json:"customerName"`
	AssignedToUserID  uuid.UUID  `json:"assignedToUserId" gorm:"type:uuid;default:uuid_generate_v4()"`
	StatusID          uuid.UUID  `json:"statusId"`
	Status            CaseStatus `gorm:"foreignKey:StatusID;references:ID"`
	DispositionMainID uuid.UUID  `json:"dispositionMainId" gorm:"type:uuid"`
	DispositionSubID  uuid.UUID  `json:"dispositionSubId" gorm:"type:uuid"`
	StartDate         time.Time  `json:"startDate" gorm:"type:date"`
	ClosedDate        time.Time  `json:"closedDate" gorm:"type:date"`
	EndDate           time.Time  `json:"endDate" gorm:"type:date"`
	ProductID         uuid.UUID  `json:"productId" gorm:"type:uuid"`
	Product           Products   `gorm:"foreignKey:ProductID;references:ID"`
	Description       string     `json:"description"`
	CreatedBy         uuid.UUID  `json:"createdBy" gorm:"type:uuid"`
	CreatedAt         time.Time  `json:"createdAt" gorm:"type:timestamp"`
	UpdatedBy         uuid.UUID  `json:"updatedBy" gorm:"type:uuid"`
	UpdatedAt         time.Time  `json:"updatedAt" gorm:"type:timestamp"`
	DeletedBy         uuid.UUID  `json:"deletedBy" gorm:"type:uuid"`
	DeletedAt         time.Time  `json:"deletedAt" gorm:"type:timestamp"`
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

type CaseTypes struct {
	ID          uuid.UUID `gorm:"primaryKey;default:uuid_generate_v4()" json:"id"`
	Group       string    `json:"group"`
	Name        string    `json:"name"`
	Description string    `json:"description" gorm:"type:text"`
}

type CaseStatus struct {
	ID          uuid.UUID `gorm:"primaryKey;default:uuid_generate_v4()" json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description" gorm:"type:text"`
}

type VerifyQuestionHistory struct {
	ID               uuid.UUID `json:"id" gorm:"type:uuid;default:uuid_generate_v4()"`
	CustomerId       string    `json:"customer_id"`
	Question         string    `json:"question" gorm:"type:text"`
	AnswerProvided   string    `json:"answer_provided" gorm:"type:text"`
	IsCorrect        bool      `json:"is_correct"`
	VeryfyByUserId   uuid.UUID `json:"verify_by_user_id" gorm:"type:uuid;default:uuid_generate_v4()"`
	VerificationDate time.Time `json:"verification_date"`
	CaseId           uuid.UUID `json:"case_id" gorm:"type:uuid;default:uuid_generate_v4()"`
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

// ##### Case Management Request #####
type CreateCaseInquiryRequest struct {
	CustomerID        string         `json:"customerId"`
	CustomerName      string         `json:"customerName"`
	CaseTypeID        uuid.UUID      `json:"caseTypeId"`
	DispositionMainID uuid.UUID      `json:"dispositionMainId" gorm:"type:uuid"`
	DispositionMains  datatypes.JSON `json:"dispositionMains" gorm:"type:jsonb"`
	DispositionSubID  uuid.UUID      `json:"dispositionSubId" gorm:"type:uuid"`
	DispositionSubs   datatypes.JSON `json:"dispositionSubs" gorm:"type:jsonb"`
	ProductID         uuid.UUID      `json:"productId" gorm:"type:uuid"`
	QueueID           uuid.UUID      `json:"queueId" gorm:"type:uuid"`
	CaseDescription   string         `json:"caseDescription" gorm:"type:text"`
	CaseNote          datatypes.JSON `json:"caseNote" gorm:"type:jsonb"`
}

type CaseFilter struct {
	Keyword     string     `form:"keyword" json:"keyword"`
	StatusID    *uint      `form:"statusId" json:"statusId"`
	PriorityID  *uint      `form:"priorityId" json:"priorityId"`
	SLADateFrom *time.Time `form:"slaDateFrom" json:"slaDateFrom"`
	SLADateTo   *time.Time `form:"slaDateTo" json:"slaDateTo"`
	Sort        string     `form:"sort" json:"sort"`
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

func (CaseTypes) TableName() string {
	return "cases_types"
}

func (CaseStatus) TableName() string {
	return "cases_status"
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
