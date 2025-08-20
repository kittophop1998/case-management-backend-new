package repository

import (
	"case-management/internal/domain/model"

	"github.com/gin-gonic/gin"
)

type CustomerRepository interface {
	CreateCustomerNote(ctx *gin.Context, note *model.CustomerNote) error
	GetAllCustomerNotes(ctx *gin.Context, customerID string, limit, offset int) ([]*model.CustomerNote, int, error)
	GetNoteTypes(ctx *gin.Context) ([]*model.NoteTypes, error)
	CountNotes(ctx *gin.Context, customerID string) (int, error)
}
