package repository

import (
	"case-management/internal/domain/model"
	"context"

	"github.com/google/uuid"
)

type CaseRepository interface {
	GetAllCase(ctx context.Context, offset, limit int, filter model.CaseFilter, category string, currID uuid.UUID) ([]*model.Cases, int, error)
	GetCaseByID(ctx context.Context, id uuid.UUID) (*model.Cases, error)
	CreateCase(ctx context.Context, c *model.Cases) (uuid.UUID, error)
	CreateCaseDispositionMains(ctx context.Context, input []*model.CaseDispositionMain) error
	CreateCaseDispositionSubs(ctx context.Context, input []*model.CaseDispositionSub) error
	GetCaseDispositionMains(ctx context.Context, caseID uuid.UUID) ([]*model.CaseDispositionMain, error)
	GetCaseDispositionSubs(ctx context.Context, caseID uuid.UUID) ([]*model.CaseDispositionSub, error)
	UpdateCaseDetail(ctx context.Context, caseToSave *model.Cases) error
	GetAllDisposition(ctx context.Context) ([]model.DispositionItem, error)
	GetCaseNotes(ctx context.Context, caseID uuid.UUID) ([]*model.CaseNotes, error)
	AddCaseNote(ctx context.Context, note *model.CaseNotes) (uuid.UUID, error)
	GenCaseCode(ctx context.Context) (string, error)
	LoadCaseStatus(ctx context.Context) (map[string]uuid.UUID, error)
	LoadCaseType(ctx context.Context) (map[string]uuid.UUID, error)
	GetCaseTypeByID(ctx context.Context, id uuid.UUID) (*model.CaseTypes, error)
}
