package usecase

import (
	"case-management/internal/domain/model"
	"case-management/internal/domain/repository"

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

func (l *LogUseCase) GetAllApiLogs(ctx *gin.Context) ([]*model.ApiLogs, error) {
	logs, err := l.repo.GetAllApiLogs(ctx)
	return logs, err
}
