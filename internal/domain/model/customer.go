package model

import (
	"time"

	"github.com/google/uuid"
)

type CustomerNote struct {
	ID          uuid.UUID `json:"id" gorm:"type:uuid;default:uuid_generate_v4()"`
	CustomerID  uuid.UUID `json:"customerId"`
	NoteTypesId uuid.UUID `json:"noteTypeId"`
	Note        string    `json:"note"`
	CreatedAt   time.Time `json:"createdAt" gorm:"autoCreateTime"`
}
