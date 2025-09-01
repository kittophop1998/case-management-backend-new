package repository

import (
	"case-management/internal/domain/model"

	"github.com/gin-gonic/gin"
)

type AuditLogRepository interface {
	LogAction(ctx *gin.Context, entry model.AuditLogs)
	Shutdown()
}
