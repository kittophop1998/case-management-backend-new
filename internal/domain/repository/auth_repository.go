package repository

import (
	"case-management/internal/domain/model"

	"github.com/gin-gonic/gin"
)

type AuthRepository interface {
	SaveAccessLog(c *gin.Context, accessLog *model.AccessLogs) error
}
