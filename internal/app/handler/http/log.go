package http

import (
	"case-management/internal/app/usecase"

	"github.com/gin-gonic/gin"
)

type LogHandler struct {
	UseCase usecase.LogUseCase
}

func (h *LogHandler) GetAllApiLogs(ctx *gin.Context) {
	logs, err := h.UseCase.GetAllApiLogs(ctx)
	if err != nil {
		ctx.JSON(500, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(200, logs)
}
