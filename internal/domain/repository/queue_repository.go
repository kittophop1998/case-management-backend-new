package repository

import (
	"case-management/internal/domain/model"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type QueueRepository interface {
	GetQueues(ctx *gin.Context, offset int, limit int, queueName string) ([]*model.Queues, int, error)
	GetQueueByID(ctx *gin.Context, id uuid.UUID) (*model.Queues, error)
	CreateQueue(ctx *gin.Context, queue *model.Queues) (uuid.UUID, error)
	CreateQueueUser(ctx *gin.Context, queueUsers []*model.QueueUsers) error
	UpdateQueue(ctx *gin.Context, input *model.Queues) error
	UpdateQueueUser(ctx *gin.Context, queueID uuid.UUID, input []*model.QueueUsers, usersDel []string) error
	IsExistingQueue(ctx *gin.Context, queueName string) bool
}
