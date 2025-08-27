package http

import (
	"case-management/infrastructure/lib"
	"case-management/internal/app/usecase"
	"case-management/internal/domain/model"
	"case-management/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

type QueueHandler struct {
	UserCase usecase.QueueUsecase
}

func (h *QueueHandler) GetAllQueues(ctx *gin.Context) {
	p := utils.GetPagination(ctx)

	queueName := ctx.Query("queueName")
	queues, total, err := h.UserCase.GetQueues(ctx, p.Page, p.Limit, queueName)
	if err != nil {
		lib.HandleError(ctx, err)
		return
	}

	lib.HandlePaginatedResponse(ctx, p.Page, p.Limit, total, queues)
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

func (h *QueueHandler) CreateQueue(ctx *gin.Context) {
	userId := ctx.GetString("userId")

	var input *model.CreateQueueRequest
	if err := ctx.ShouldBindJSON(&input); err != nil {
		lib.HandleError(ctx, lib.BadRequest.WithDetails(err.Error()))
		return
	}

	queueId, err := h.UserCase.CreateQueue(ctx, userId, input)
	if err != nil {
		lib.HandleError(ctx, lib.CannotCreate.WithDetails(err.Error()))
		return
	}

	lib.HandleResponse(ctx, http.StatusCreated, gin.H{"queueId": queueId})
}
