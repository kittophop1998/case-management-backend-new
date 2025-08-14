package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"
)

// Case Entity
type Cases struct {
	Model
	Title               string         `json:"title"`
	CustomerId          string         `json:"customerId"`
	DispositionMainId   uuid.UUID      `json:"dispositionMainId" gorm:"type:uuid"`
	DispositionSubId    uuid.UUID      `json:"dispositionSubId" gorm:"type:uuid"`
	CaseTypeId          uuid.UUID      `json:"caseTypeId" gorm:"type:uuid;default:uuid_generate_v4()"`
	CreditCardAccountId string         `json:"creditCardAccountId"`
	LoanAccountId       string         `json:"loanAccountId"`
	AssignedToUserId    uuid.UUID      `json:"assignedToUserId" gorm:"type:uuid;default:uuid_generate_v4()"`
	StatusId            uint           `json:"statusId"`
	PriorityId          uint           `json:"priorityId"`
	StartDate           time.Time      `json:"startDate" gorm:"type:date"`
	EndDate             time.Time      `json:"endDate" gorm:"type:date"`
	CaseNote            datatypes.JSON `json:"caseNote" gorm:"type:jsonb"`
	Resolution          string         `json:"resolution" gorm:"type:text"`
	CreatedBy           string         `json:"createdBy"`
	SLADate             time.Time      `json:"slaDate"`
}

// Create Case Request Body
type CreateCaseRequest struct {
	Title             string         `json:"title"`
	CustomerId        string         `json:"customerId"`
	CaseTypeId        uuid.UUID      `json:"caseTypeId"`
	DispositionMainId uuid.UUID      `json:"dispositionMainId" gorm:"type:uuid"`
	DispositionMains  datatypes.JSON `json:"dispositionMains" gorm:"type:jsonb"`
	DispositionSubId  uuid.UUID      `json:"dispositionSubId" gorm:"type:uuid"`
	DispositionSubs   datatypes.JSON `json:"dispositionSubs" gorm:"type:jsonb"`
	CaseDescription   string         `json:"caseDescription" gorm:"type:text"`
	CaseNote          datatypes.JSON `json:"caseNote" gorm:"type:jsonb"`
}

// Case Disposition Entity
type CaseDispositionMain struct {
	Model
	CaseId            uuid.UUID `json:"caseId" gorm:"type:uuid"`
	DispositionMainId uuid.UUID `json:"dispositionMainId" gorm:"type:uuid"`
}

// Case Disposition Sub Entity
type CaseDispositionSub struct {
	Model
	CaseId           uuid.UUID `json:"caseId" gorm:"type:uuid"`
	DispositionSubId uuid.UUID `json:"dispositionSubId" gorm:"type:uuid"`
}

type CaseFilter struct {
	Keyword     string     `form:"keyword" json:"keyword"`
	StatusID    *uint      `form:"statusId" json:"statusId"`
	PriorityID  *uint      `form:"priorityId" json:"priorityId"`
	SLADateFrom *time.Time `form:"slaDateFrom" json:"slaDateFrom"`
	SLADateTo   *time.Time `form:"slaDateTo" json:"slaDateTo"`
	Sort        string     `form:"sort" json:"sort"`
}

type NoteTypes struct {
	ID          uuid.UUID `gorm:"primaryKey;default:uuid_generate_v4()" json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description" gorm:"type:text"`
}

type CaseTypes struct {
	ID          uuid.UUID `gorm:"primaryKey;default:uuid_generate_v4()" json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description" gorm:"type:text"`
}

type CaseStatus struct {
	ID             uuid.UUID `gorm:"primaryKey;default:uuid_generate_v4()" json:"id"`
	Name           string    `json:"name"`
	Description    string    `json:"description" gorm:"type:text"`
	IsClosedStatus bool      `json:"is_closed_status"`
}

type CasePriorities struct {
	ID          uuid.UUID `gorm:"primaryKey;default:uuid_generate_v4()" json:"id"`
	Name        string    `json:"name"`
	OrderNumber uint      `json:"order_number"`
}

type CaseNotes struct {
	ID          uuid.UUID `gorm:"primaryKey;default:uuid_generate_v4()" json:"id"`
	CaseId      uuid.UUID `json:"case_id" gorm:"type:uuid;default:uuid_generate_v4()"`
	UserId      uuid.UUID `json:"user_id" gorm:"type:uuid;default:uuid_generate_v4()"`
	NoteTypesId uuid.UUID `json:"note_types_id"`
	Content     string    `json:"content" gorm:"type:text"`
	CreatedAt   time.Time `json:"created_at"`
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

func (CasePriorities) TableName() string {
	return "cases_priorities"
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
