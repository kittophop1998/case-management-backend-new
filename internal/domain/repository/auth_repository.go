package repository

import (
	"case-management/internal/domain/model"

	"github.com/gin-gonic/gin"
)

type AuthRepository interface {
	SaveAccessLog(ctx *gin.Context, accessLog *model.AccessLogs) error
}
