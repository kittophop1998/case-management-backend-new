package repository

import (
	"case-management/internal/domain/model"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type QueueRepository interface {
	GetQueues(ctx *gin.Context, offset int, limit int, queueName string) ([]*model.Queues, int, error)
	GetQueueByID(ctx *gin.Context, id uuid.UUID) (*model.Queues, error)
}
