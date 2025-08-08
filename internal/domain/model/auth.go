package model

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginResponse struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
}

type AccessLogs struct {
	ID            uuid.UUID       `json:"id" gorm:"type:uuid;default:uuid_generate_v4()"`
	UserID        uuid.UUID       `json:"userId" gorm:"type:uuid;default:uuid_generate_v4()"`
	Action        string          `json:"action"`
	IPAddress     string          `json:"ipAddress"`
	UserAgent     string          `json:"userAgent" gorm:"type:text"`
	Details       json.RawMessage `gorm:"type:jsonb" json:"details"`
	CreatedAt     time.Time       `json:"createdAt"`
	Username      string          `gorm:"type:varchar(20)" json:"username"`
	LogonDatetime time.Time       `gorm:"type:timestamp" json:"logonDatetime"`
	LogonResult   string          `gorm:"type:varchar(10)" json:"logonResult"`
}
