package repository

import (
	"case-management/internal/domain/model"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/datatypes"
)

type CaseRepository interface {
	CreateCase(ctx *gin.Context, c *model.CreateCaseRequest) (uuid.UUID, error)
	GetAllCase(ctx *gin.Context, offset, limit int) ([]*model.Cases, int, error)
	GetCaseByID(ctx *gin.Context, id uuid.UUID) (*model.Cases, error)
	AddInitialDescription(ctx *gin.Context, caseID uuid.UUID, newDescription string) error
	CreateCaseDispositionMains(ctx *gin.Context, data datatypes.JSON) error
	CreateCaseDispositionSubs(ctx *gin.Context, data datatypes.JSON) error
	GetAllDisposition(ctx *gin.Context, limit, offset int) ([]model.DispositionMain, int, error)
	// GetNoteTypeByID(ctx *gin.Context, noteTypeID uuid.UUID) (*model.NoteTypes, error)
}
