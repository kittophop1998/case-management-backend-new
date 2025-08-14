package model

import (
	"github.com/google/uuid"
)

type CustomerNote struct {
	Model
	ID          uuid.UUID `json:"id" gorm:"type:uuid;default:uuid_generate_v4()"`
	CustomerID  string    `json:"customerId"`
	NoteTypesId uuid.UUID `json:"noteTypeId"`
	Note        string    `json:"note"`
}
