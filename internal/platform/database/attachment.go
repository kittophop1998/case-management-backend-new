package database

import (
	"case-management/internal/domain/model"
	"context"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type AttachmentsPg struct {
	db *gorm.DB
}

func NewAttachmentsPg(db *gorm.DB) *AttachmentsPg {
	return &AttachmentsPg{db: db}
}

func (a AttachmentsPg) UploadAttachment(ctx context.Context, caseID uuid.UUID, file model.Attachment) (uuid.UUID, error) {
	tx := a.db.WithContext(ctx).Create(&file)
	return file.ID, tx.Error
}
