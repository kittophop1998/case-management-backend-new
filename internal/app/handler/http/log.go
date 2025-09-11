package http

import (
	"case-management/infrastructure/lib"
	"case-management/internal/app/usecase"
	"case-management/internal/domain/model"
	"case-management/utils"
	"log"

	"github.com/gin-gonic/gin"
)

type LogHandler struct {
	UseCase usecase.LogUseCase
}

func (h *LogHandler) GetAllApiLogs(ctx *gin.Context) {
	var query model.APILogQueryParams
	p := utils.GetPagination(ctx)

	if err := ctx.ShouldBindQuery(&query); err != nil {
		log.Println("[Handler] Failed to bind query params:", err)
		lib.HandleError(ctx, lib.BadRequest.WithDetails("Invalid query parameters"))
		return
	}

	log.Printf("[Handler] Parsed query: %+v\n", query)
	log.Printf("[Handler] Pagination: page=%d, limit=%d\n", p.Page, p.Limit)

	logs, total, err := h.UseCase.GetAllApiLogs(ctx, p.Page, p.Limit, &query)
	if err != nil {
		lib.HandleError(ctx, lib.InternalServer.WithDetails(err.Error()))
		return
	}

	lib.HandlePaginatedResponse(ctx, p.Page, p.Limit, total, logs)
}
