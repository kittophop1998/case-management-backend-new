package repository

import (
	"case-management/internal/domain/model"
	"context"

	"github.com/google/uuid"
)

// Attachment
type AttachmentsRepository interface {
	UploadAttachment(ctx context.Context, caseID uuid.UUID, file model.Attachment) (uuid.UUID, error)
}
