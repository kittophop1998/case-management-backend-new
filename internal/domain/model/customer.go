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
}

// ##### Master Data for Customer #####
type NoteTypes struct {
	Model
	Name        string `json:"name"`
	Description string `json:"description" gorm:"type:text"`
}
