package model

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"
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

type AccessLogs struct {
	ID            uuid.UUID       `json:"id" gorm:"type:uuid;default:uuid_generate_v4()"`
	UserID        uuid.UUID       `json:"userId" gorm:"type:uuid;default:uuid_generate_v4()"`
	Action        string          `json:"action"`
	Details       json.RawMessage `gorm:"type:jsonb" json:"details"`
	CreatedAt     time.Time       `json:"createdAt"`
	Username      string          `gorm:"type:varchar(20)" json:"username"`
	LogonDatetime time.Time       `gorm:"type:timestamp" json:"logonDatetime"`
	LoginSuccess  *bool           `json:"loginSuccess"`
}

type AuditLogs struct {
	ID        uuid.UUID      `json:"id" gorm:"type:uuid;default:uuid_generate_v4()"`
	UserID    uuid.UUID      `json:"userId" gorm:"type:uuid;default:uuid_generate_v4()"`
	Action    string         `json:"action"`
	Module    string         `json:"module"`
	Metadata  datatypes.JSON `gorm:"type:jsonb" json:"metadata"`
	CreatedAt time.Time      `json:"createdAt"`
}
