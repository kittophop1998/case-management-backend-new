package database

import (
	"case-management/internal/domain/model"
	"fmt"

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

func (repo *QueuePg) GetQueues(
	ctx *gin.Context,
	offset int,
	limit int,
	queueName string,
) ([]*model.Queues, int, error) {
	var queues []*model.Queues

	baseQuery := repo.db.WithContext(ctx).Model(&model.Queues{})

	// Filter by name
	if queueName != "" {
		baseQuery = baseQuery.Where("name ILIKE ?", "%"+queueName+"%")
	}

	// Count total rows
	var count int64
	if err := baseQuery.Count(&count).Error; err != nil {
		return nil, 0, err
	}

	// ตรวจสอบ offset ไม่เกิน count
	if int64(offset) >= count {
		offset = 0
	}

	// Query data
	query := baseQuery.Session(&gorm.Session{}) // clone query
	if err := query.Limit(limit).Offset(offset).Debug().Find(&queues).Error; err != nil {
		return nil, 0, err
	}

	// บังคับ default empty slice
	if queues == nil {
		queues = []*model.Queues{}
	}

	fmt.Println("Queues found:", len(queues), "Total count:", count)

	return queues, int(count), nil
}

func (repo *QueuePg) GetQueueByID(ctx *gin.Context, id uuid.UUID) (*model.Queues, error) {
	var queue *model.Queues
	if err := repo.db.WithContext(ctx).Where(&queue, model.Queues{ID: id}).First(&queue).Error; err != nil {
		return nil, err
	}
	return queue, nil
}
