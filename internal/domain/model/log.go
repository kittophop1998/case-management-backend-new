package model

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"
)

type ApiLogs struct {
	ID           string    `gorm:"primaryKey" json:"id"`
	RequestID    string    `json:"request_id"`
	ServiceName  string    `json:"service_name"`
	Endpoint     string    `json:"endpoint"`
	ReqDatetime  time.Time `json:"req_datetime"`
	ReqHeader    string    `json:"req_header"`
	ReqMessage   string    `json:"req_message"`
	RespDatetime time.Time `json:"resp_datetime"`
	RespHeader   string    `json:"resp_header"`
	RespMessage  string    `json:"resp_message"`
	StatusCode   int       `json:"status_code"`
	TimeUsage    int       `json:"time_usage"`
	CreatedAt    time.Time `json:"created_at"`
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

type APILogQueryParams struct {
	RequestID    string `form:"request_id"`
	ServiceName  string `form:"service_name"`
	Endpoint     string `form:"endpoint"`
	ReqDatetime  string `form:"req_datetime"`
	ReqHeader    string `form:"req_header"`
	ReqMessage   string `form:"req_message"`
	RespDatetime string `form:"resp_datetime"`
	RespHeader   string `form:"resp_header"`
	RespMessage  string `form:"resp_message"`
	StatusCode   int    `form:"status_code"`
	TimeUsage    int    `form:"time_usage"`
	SortingField string `form:"sorting_field"`
	SortingOrder string `form:"sorting_order"`
}
