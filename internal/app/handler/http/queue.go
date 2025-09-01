package http

import (
	"case-management/infrastructure/lib"
	"case-management/internal/app/usecase"
	"case-management/internal/domain/model"
	"case-management/utils"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
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
	idStr := ctx.Param("id")
	queueID, err := uuid.Parse(idStr)
	if err != nil {
		lib.HandleError(ctx, lib.BadRequest.WithDetails("Invalid queue ID"))
		return
	}

	queue, err := h.UserCase.GetQueueByID(ctx, queueID)
	if err != nil {
		lib.HandleError(ctx, err)
		return
	}

	lib.HandleResponse(ctx, http.StatusOK, queue)
}

func (h *QueueHandler) CreateQueue(ctx *gin.Context) {
	userId := ctx.GetString("userId")
	createdByID, err := uuid.Parse(userId)
	if err != nil {
		lib.HandleError(ctx, lib.BadRequest.WithDetails("Invalid user ID"))
		return
	}

	var input *model.CreateQueueRequest
	if err := ctx.ShouldBindJSON(&input); err != nil {
		lib.HandleError(ctx, lib.BadRequest.WithDetails(err.Error()))
		return
	}

	queueId, err := h.UserCase.CreateQueue(ctx, createdByID, input)
	if err != nil {
		lib.HandleError(ctx, lib.CannotCreate.WithDetails(err.Error()))
		return
	}

	lib.HandleResponse(ctx, http.StatusCreated, gin.H{"queueId": queueId})
}

func (h QueueHandler) AddUserInQueue(ctx *gin.Context) {
	userId := ctx.GetString("userId")
	createdByID, err := uuid.Parse(userId)
	if err != nil {
		lib.HandleError(ctx, lib.BadRequest.WithDetails("Invalid user ID"))
		return
	}

	idStr := ctx.Param("id")
	queueID, err := uuid.Parse(idStr)
	if err != nil {
		lib.HandleError(ctx, lib.BadRequest.WithDetails("Invalid queue ID"))
		return
	}

	var users model.UserManageInQueue
	if err := ctx.ShouldBindJSON(&users); err != nil {
		lib.HandleError(ctx, lib.BadRequest.WithDetails("1"+err.Error()))
		return
	}

	if err := h.UserCase.AddUserInQueue(ctx, createdByID, queueID, users); err != nil {
		lib.HandleError(ctx, lib.CannotUpdate.WithDetails(err.Error()))
		return
	}

	lib.HandleResponse(ctx, http.StatusOK, gin.H{"message": "Add user to queue success"})
}

func (h *QueueHandler) UpdateQueueByID(ctx *gin.Context) {
	userId := ctx.GetString("userId")
	updatedByID, err := uuid.Parse(userId)
	if err != nil {
		lib.HandleError(ctx, lib.BadRequest.WithDetails("Invalid user ID"))
		return
	}

	idStr := ctx.Param("id")
	queueID, err := uuid.Parse(idStr)
	if err != nil {
		lib.HandleError(ctx, lib.BadRequest.WithDetails("Invalid queue ID"))
		return
	}

	var queueUpdate model.Queues
	if err := ctx.ShouldBindJSON(&queueUpdate); err != nil {
		lib.HandleError(ctx, lib.BadRequest.WithDetails(err.Error()))
		return
	}

	if err := h.UserCase.UpdateQueueByID(ctx, updatedByID, queueID, &queueUpdate); err != nil {
		lib.HandleError(ctx, lib.CannotUpdate.WithDetails(err.Error()))
		return
	}

	lib.HandleResponse(ctx, http.StatusOK, gin.H{"message": "Update queue success"})
}

func (h *QueueHandler) DeleteUsersInQueue(ctx *gin.Context) {
	userId := ctx.GetString("userId")
	createdByID, err := uuid.Parse(userId)
	if err != nil {
		lib.HandleError(ctx, lib.BadRequest.WithDetails("Invalid user ID"))
		return
	}

	idStr := ctx.Param("id")
	queueID, err := uuid.Parse(idStr)
	if err != nil {
		lib.HandleError(ctx, lib.BadRequest.WithDetails("Invalid queue ID"))
		return
	}

	var users model.UserManageInQueue
	if err := ctx.ShouldBindJSON(&users); err != nil {
		lib.HandleError(ctx, lib.BadRequest.WithDetails(err.Error()))
		return
	}

	if err := h.UserCase.DeleteUsersInQueue(ctx, createdByID, queueID, users); err != nil {
		lib.HandleError(ctx, err)
		return
	}

	lib.HandleResponse(ctx, http.StatusOK, gin.H{"message": "Delete users from queue success"})
}
