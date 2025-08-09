package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"
)

type Cases struct {
	Model
	Title               string         `json:"title"`
	CustomerId          string         `json:"customer_id"`
	CreditCardAccountId string         `json:"credit_card_account_id"`
	LoanAccountId       string         `json:"loan_account_id"`
	AssignedToUserId    uuid.UUID      `json:"assigned_to_user_id" gorm:"type:uuid;default:uuid_generate_v4()"`
	StatusId            uint           `json:"status_id"`
	PriorityId          uint           `json:"priority_id"`
	StartDate           time.Time      `json:"start_date" gorm:"type:date"`
	EndDate             time.Time      `json:"end_date" gorm:"type:date"`
	InitialDescriptions datatypes.JSON `json:"initial_descriptions" gorm:"type:jsonb"`
	Resolution          string         `json:"resolution" gorm:"type:text"`
	CreatedBy           string         `json:"created_by"`
	SLADate             time.Time      `json:"sla_date"`
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

type CustomerNote struct {
	ID          uuid.UUID `json:"id" gorm:"type:uuid;default:uuid_generate_v4()"`
	CustomerID  uuid.UUID `json:"customerId"`
	NoteTypesId uuid.UUID `json:"noteTypeId"`
	Note        string    `json:"note"`
	CreatedAt   time.Time `json:"createdAt" gorm:"autoCreateTime"`
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
