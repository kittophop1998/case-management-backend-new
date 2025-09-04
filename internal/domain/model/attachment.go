package model

import (
	"time"

	"github.com/google/uuid"
)

type Attachment struct {
	ID               uuid.UUID `json:"id" gorm:"type:uuid;default:uuid_generate_v4()"`
	CaseId           uuid.UUID `json:"case_id" gorm:"type:uuid;default:uuid_generate_v4()"`
	FileName         string    `json:"file_name"`
	FilePath         string    `json:"file_path"`
	FileType         string    `json:"file_type"`
	FileSizeBytes    uint64    `json:"file_size_bytes" gorm:"type:bigint"`
	UploadedByUserId uuid.UUID `json:"uploaded_by_user_id" gorm:"type:uuid;default:uuid_generate_v4()"`
	UploadedAt       time.Time `json:"uploaded_at"`
}

func (Attachment) TableName() string {
	return "attachments"
}
