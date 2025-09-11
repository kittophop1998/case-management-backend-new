package repository

import (
	"case-management/internal/domain/model"
	"context"

	"github.com/gin-gonic/gin"
)

type LogRepository interface {
	SaveApiLog(log *model.ApiLogs) error
	GetAllApiLogs(ctx context.Context, limit, offset int, q *model.APILogQueryParams) ([]*model.ApiLogs, int, error)
	SaveLoginEvent(ctx *gin.Context, accessLog *model.AccessLogs) error
}
