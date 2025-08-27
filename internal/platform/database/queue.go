package database

import (
	"case-management/internal/domain/model"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type QueuePg struct {
	db *gorm.DB
}

func NewQueuePg(db *gorm.DB) *QueuePg {
	return &QueuePg{db: db}
}

func (repo *QueuePg) GetQueues(ctx *gin.Context) ([]*model.Queues, error) {
	var queues []*model.Queues
	if err := repo.db.WithContext(ctx).Find(&queues).Error; err != nil {
		return nil, err
	}
	return queues, nil
}

func (repo *QueuePg) GetQueueByID(ctx *gin.Context, id uuid.UUID) (*model.Queues, error) {
	var queue *model.Queues
	if err := repo.db.WithContext(ctx).Where(&queue, model.Queues{ID: id}).First(&queue).Error; err != nil {
		return nil, err
	}
	return queue, nil
}
