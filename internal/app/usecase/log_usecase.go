package usecase

import (
	"case-management/internal/domain/model"
	"case-management/internal/domain/repository"
	"context"

	"github.com/gin-gonic/gin"
)

type LogUseCase struct {
	repo repository.LogRepository
}

func NewLogUseCase(repo repository.LogRepository) *LogUseCase {
	return &LogUseCase{repo: repo}
}

func (l *LogUseCase) SaveApiLog(log *model.ApiLogs) error {
	return l.repo.SaveApiLog(log)
}

func (l *LogUseCase) GetAllApiLogs(ctx context.Context, page, limit int, q *model.APILogQueryParams) ([]*model.ApiLogs, int, error) {
	offset := (page - 1) * limit
	logs, total, err := l.repo.GetAllApiLogs(ctx, limit, offset, q)
	return logs, total, err
}

func (repo *LogUseCase) SaveLoginEvent(ctx *gin.Context, accessLog *model.AccessLogs) error {
	return repo.repo.SaveLoginEvent(ctx, accessLog)
}
