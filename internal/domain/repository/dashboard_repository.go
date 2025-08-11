package repository

import (
	"case-management/internal/domain/model"
	"context"
)

type DashboardRepository interface {
	FindByAeonID(ctx context.Context, aeonID string) (*model.GetCustomerInfoResponse, error)
}
