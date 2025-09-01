package repository

import (
	"case-management/internal/domain/model"

	"github.com/gin-gonic/gin"
)

type LogRepository interface {
	SaveApiLog(log *model.ApiLogs) error
	GetAllApiLogs(ctx *gin.Context) ([]*model.ApiLogs, error)
	SaveLoginEvent(ctx *gin.Context, accessLog *model.AccessLogs) error
}
