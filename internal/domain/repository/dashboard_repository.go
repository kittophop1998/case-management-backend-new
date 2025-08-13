package repository

import (
	"case-management/internal/domain/model"
	"context"
)

type DashboardRepository interface {
	GetCustInfoByAeonID(ctx context.Context, aeonID string) (*model.GetCustInfoResponse, error)
	GetCustProfileByAeonID(ctx context.Context, aeonID string) (*model.GetCustProfileResponse, error)
	GetCustSegmentByAeonID(ctx context.Context, aeonID string) (*model.GetCustSegmentResponse, error)
	GetCustSuggestionByAeonID(ctx context.Context, aeonID string) (*model.GetCustSuggestionResponse, error)
}
