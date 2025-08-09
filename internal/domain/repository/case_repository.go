package repository

import (
	"case-management/internal/domain/model"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type CaseRepository interface {
	CreateCase(ctx *gin.Context, c *model.Cases) (uuid.UUID, error)
	GetAllCase(ctx *gin.Context, limit, offset int, filter model.CaseFilter) ([]*model.Cases, error)
	CountWithFilter(ctx *gin.Context, filter model.CaseFilter) (int, error)
	CreateNoteType(ctx *gin.Context, note model.NoteTypes) (*model.NoteTypes, error)
	GetCaseByID(ctx *gin.Context, id uuid.UUID) (*model.Cases, error)
	AddInitialDescription(ctx *gin.Context, caseID uuid.UUID, newDescription string) error
	GetNoteTypeByID(ctx *gin.Context, noteTypeID uuid.UUID) (*model.NoteTypes, error)
}
