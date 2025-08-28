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

func (repo *QueuePg) GetQueues(
	ctx *gin.Context,
	offset int,
	limit int,
	queueName string,
) ([]*model.Queues, int, error) {
	var queues []*model.Queues

	// Count query
	countQuery := repo.db.WithContext(ctx).Model(&model.Queues{})
	if queueName != "" {
		countQuery = countQuery.Where("name ILIKE ?", "%"+queueName+"%")
	}

	var total int64
	if err := countQuery.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	dataQuery := repo.db.WithContext(ctx).Model(&model.Queues{})
	if queueName != "" {
		dataQuery = dataQuery.Where("name ILIKE ?", "%"+queueName+"%")
	}

	if err := dataQuery.Limit(limit).Offset(offset).Find(&queues).Error; err != nil {
		return nil, 0, err
	}

	if queues == nil {
		queues = []*model.Queues{}
	}

	return queues, int(total), nil
}

func (repo *QueuePg) GetQueueByID(ctx *gin.Context, id uuid.UUID) (*model.Queues, error) {
	var queue *model.Queues
	if err := repo.db.WithContext(ctx).Where(&queue, model.Queues{ID: id}).First(&queue).Error; err != nil {
		return nil, err
	}
	return queue, nil
}

func (repo *QueuePg) CreateQueue(ctx *gin.Context, queue *model.Queues) (uuid.UUID, error) {
	if err := repo.db.WithContext(ctx).FirstOrCreate(queue).Error; err != nil {
		return uuid.Nil, err
	}
	return queue.ID, nil
}

func (repo *QueuePg) CreateQueueUser(ctx *gin.Context, queueUsers []*model.QueueUsers) error {
	if err := repo.db.WithContext(ctx).Create(queueUsers).Error; err != nil {
		return err
	}
	return nil
}

func (repo *QueuePg) UpdateQueue(ctx *gin.Context, queue *model.Queues) error {
	if err := repo.db.WithContext(ctx).
		Where("id = ?", queue.ID).
		Updates(queue).Error; err != nil {
		return err
	}
	return nil
}

func (repo *QueuePg) UpdateQueueUser(ctx *gin.Context, queueID uuid.UUID, usersAdd []*model.QueueUsers, usersDel []string) error {
	return repo.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if len(usersDel) > 0 {
			if err := tx.Where("queue_id = ? AND user_id IN ?", queueID, usersDel).
				Delete(&model.QueueUsers{}).Error; err != nil {
				return err
			}
		}

		if len(usersAdd) > 0 {
			if err := tx.Create(usersAdd).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

func (repo *QueuePg) IsExistingQueue(ctx *gin.Context, queueName string) bool {
	var count int64
	if err := repo.db.WithContext(ctx).
		Model(&model.Queues{}).
		Where("name = ?", queueName).
		Count(&count).Error; err != nil {
		return false
	}
	return count > 0

}
