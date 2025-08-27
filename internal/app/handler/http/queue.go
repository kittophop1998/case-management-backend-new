package http

import (
	"case-management/infrastructure/lib"
	"case-management/internal/app/usecase"
	"net/http"

	"github.com/gin-gonic/gin"
)

type QueueHandler struct {
	UserCase usecase.QueueUsecase
}

func (h *QueueHandler) GetAllQueues(ctx *gin.Context) {
	limit, err := getLimit(ctx)
	if err != nil {
		lib.HandleError(ctx, lib.BadRequest.WithDetails(err.Error()))
		return
	}

	page, err := getPage(ctx)
	if err != nil {
		lib.HandleError(ctx, lib.BadRequest.WithDetails(err.Error()))
		return
	}

	queueName := ctx.Query("queueName")
	queues, total, err := h.UserCase.GetQueues(ctx, page, limit, queueName)
	if err != nil {
		lib.HandleError(ctx, err)
		return
	}

	lib.HandlePaginatedResponse(ctx, page, limit, total, queues)
}

func (h *QueueHandler) GetQueueByID(ctx *gin.Context) {
	id := ctx.Param("id")
	queue, err := h.UserCase.GetQueueByID(ctx, id)
	if err != nil {
		lib.HandleError(ctx, err)
		return
	}

	lib.HandleResponse(ctx, http.StatusOK, queue)
}
