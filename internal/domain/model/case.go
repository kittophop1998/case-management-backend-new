package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"
)

type Cases struct {
	Model
	Title               string         `json:"title"`
	CustomerId          string         `json:"customerId"`
	DispositionMainId   datatypes.JSON `json:"dispositionMainId" gorm:"type:jsonb"`
	DispositionSubId    datatypes.JSON `json:"dispositionSubId" gorm:"type:jsonb"`
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
	Title           string         `json:"title"`
	CustomerId      string         `json:"customerId"`
	CaseTypeId      uuid.UUID      `json:"caseTypeId"`
	DispositionMain datatypes.JSON `json:"dispositionMain" gorm:"type:jsonb"`
	DispositionSub  datatypes.JSON `json:"dispositionSub" gorm:"type:jsonb"`
	CaseDescription string         `json:"caseDescription" gorm:"type:text"`
	CaseNote        datatypes.JSON `json:"caseNote" gorm:"type:jsonb"`
}

type CaseFilter struct {
	Keyword     string     `form:"keyword" json:"keyword"`
	StatusID    *uint      `form:"status_id" json:"status_id"`
	PriorityID  *uint      `form:"priority_id" json:"priority_id"`
	SLADateFrom *time.Time `form:"sla_date_from" json:"sla_date_from"`
	SLADateTo   *time.Time `form:"sla_date_to" json:"sla_date_to"`
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
