package repository

import (
	"case-management/internal/domain/model"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type CustomerRepository interface {
	CreateCustomerNote(ctx *gin.Context, note *model.CustomerNote) error
	GetAllCustomerNotes(ctx *gin.Context, limit, offset int, customerID uuid.UUID, filter model.CustomerNoteFilter) ([]*model.CustomerNote, int, error)
	GetNoteTypes(ctx *gin.Context) ([]*model.NoteTypes, error)
	CountNotes(ctx *gin.Context, customerID uuid.UUID) (int, error)
}
