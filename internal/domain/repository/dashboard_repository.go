package repository

import (
	"case-management/internal/domain/model"
	"context"
)

type DashboardRepository interface {
	GetCustInfoByAeonID(ctx context.Context, aeonID string) (*model.GetCustomerInfoResponse, error)
	GetCustProfileByAeonID(ctx context.Context, aeonID string) (*model.GetCustomerProfileResponse, error)
}
