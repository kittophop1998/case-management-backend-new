package repository

import (
	"case-management/internal/domain/model"
	"context"
)

type DashboardRepository interface {
	SaveApiLog(ctx context.Context, log *model.ApiLogs) error
}
