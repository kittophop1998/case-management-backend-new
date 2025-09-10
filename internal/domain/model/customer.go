package model

import (
	"github.com/google/uuid"
)

type CustomerNote struct {
	Model
	ID          uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4()" json:"id"`
	CustomerID  string    `json:"customerId" binding:"required"`
	NoteTypesID uuid.UUID `json:"noteTypeId" binding:"required"`
	NoteType    NoteTypes `gorm:"foreignKey:NoteTypesID" json:"noteType"`
	Note        string    `json:"note" binding:"required"`
	CreatedBy   string    `json:"createdBy"`
}

// ##### Master Data for Customer #####
type NoteTypes struct {
	Model
	Name        string `json:"name"`
	Description string `json:"description" gorm:"type:text"`
}

// ##### Customer Response #####
type CustomerNoteResponse struct {
	ID          string `json:"id"`
	NoteType    string `json:"noteType"`
	NoteDetail  string `json:"noteDetail"`
	CreatedBy   string `json:"createdBy"`
	CreatedDate string `json:"createdDate"`
}

// ##### Customer Note Request #####
type CustomerNoteFilter struct {
	NoteTypeID *uuid.UUID `form:"noteTypeId"`
	Keyword    string     `form:"keyword"`
}
