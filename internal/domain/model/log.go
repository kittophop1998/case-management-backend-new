package model

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

type ApiLogs struct {
	ID              uuid.UUID       `json:"id" gorm:"type:uuid;default:uuid_generate_v4()"`
	UserId          uuid.UUID       `json:"user_id" gorm:"type:uuid;default:uuid_generate_v4()"`
	Endpoint        string          `json:"endpoint"`
	Method          string          `json:"method"`
	RequestPayload  json.RawMessage `gorm:"type:jsonb" json:"request_payload"`
	ResponsePayload json.RawMessage `gorm:"type:jsonb" json:"response_payload"`
	StatusCode      uint            `json:"status_code"`
	DurationMs      uint            `json:"duration_ms"`
	ErrorMessage    string          `json:"error_message" gorm:"type:text"`
	CreatedAt       time.Time       `json:"created_at"`
}
