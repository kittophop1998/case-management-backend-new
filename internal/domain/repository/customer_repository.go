package repository

import (
	"case-management/internal/domain/model"

	"github.com/gin-gonic/gin"
)

type CustomerRepository interface {
	CreateCustomerNote(ctx *gin.Context, note *model.CustomerNote) error
	GetAllCustomerNotes(ctx *gin.Context, customerID string) ([]*model.CustomerNote, error)
	GetNoteTypes(ctx *gin.Context) ([]*model.NoteTypes, error)
	CountNotes(ctx *gin.Context, customerID string) (int, error)
}
