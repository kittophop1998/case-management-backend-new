package updatecaseupdater

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type CaseUpdater interface {
	Update(ctx *gin.Context, caseID uuid.UUID, data map[string]interface{}) error
}
