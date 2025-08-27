package repository

import (
	"case-management/internal/domain/model"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type QueueRepository interface {
	GetQueues(ctx *gin.Context) ([]*model.Queues, error)
	GetQueueByID(ctx *gin.Context, id uuid.UUID) (*model.Queues, error)
}
