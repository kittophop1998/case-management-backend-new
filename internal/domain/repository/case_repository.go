package repository

import (
	"case-management/internal/domain/model"
	"context"

	"github.com/google/uuid"
	"gorm.io/datatypes"
)

type CaseRepository interface {
	GetAllCase(ctx context.Context, offset, limit int, filter model.CaseFilter, category string, currID uuid.UUID) ([]*model.Cases, int, error)
	GetCaseByID(ctx context.Context, id uuid.UUID) (*model.Cases, error)
	CreateCase(ctx context.Context, c *model.Cases) (uuid.UUID, error)
	CreateCaseDispositionMains(ctx context.Context, data datatypes.JSON) error
	CreateCaseDispositionSubs(ctx context.Context, data datatypes.JSON) error
	UpdateCaseDetail(ctx context.Context, caseToSave *model.Cases) error
	GetAllDisposition(ctx context.Context) ([]model.DispositionItem, error)
	AddCaseNote(ctx context.Context, note *model.CaseNotes) (uuid.UUID, error)
	AddInitialDescription(ctx context.Context, caseID uuid.UUID, newDescription string) error
	GenCaseCode(ctx context.Context) (string, error)
	LoadCaseStatus(ctx context.Context) (map[string]uuid.UUID, error)
	LoadCaseType(ctx context.Context) (map[string]uuid.UUID, error)
}
