package repository

import (
	"case-management/internal/domain/model"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/datatypes"
)

type CaseRepository interface {
	CreateCaseInquiry(ctx *gin.Context, c *model.Cases) (uuid.UUID, error)
	GetAllCase(ctx *gin.Context, offset, limit int, category string, currID uuid.UUID) ([]*model.Cases, int, error)
	GetCaseByID(ctx *gin.Context, id uuid.UUID) (*model.Cases, error)
	AddInitialDescription(ctx *gin.Context, caseID uuid.UUID, newDescription string) error
	CreateCaseDispositionMains(ctx *gin.Context, data datatypes.JSON) error
	CreateCaseDispositionSubs(ctx *gin.Context, data datatypes.JSON) error
	GetAllDisposition(ctx *gin.Context) ([]model.DispositionItem, error)
	LoadCaseStatus(ctx *gin.Context) (map[string]uuid.UUID, error)
	LoadCaseType(ctx *gin.Context) (map[string]uuid.UUID, error)
	AddCaseNote(ctx *gin.Context, note *model.CaseNotes) (uuid.UUID, error)
}
