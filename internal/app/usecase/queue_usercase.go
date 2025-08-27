package usecase

import (
	"case-management/internal/domain/model"
	"case-management/internal/domain/repository"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type QueueUsercase struct {
	repo repository.QueueRepository
}

func NewQueueUsercase(repo repository.QueueRepository) *QueueUsercase {
	return &QueueUsercase{repo: repo}
}

func (u QueueUsercase) GetQueues(ctx *gin.Context) ([]*model.GetQueuesResponse, error) {
	var queues []*model.GetQueuesResponse
	queuesRepo, err := u.repo.GetQueues(ctx)
	if err != nil {
		return nil, err
	}

	for _, queue := range queuesRepo {
		queues = append(queues, &model.GetQueuesResponse{
			QueueID:          queue.ID,
			QueueName:        queue.Name,
			QueueDescription: queue.Description,
			CreatedAt:        queue.CreatedAt,
			CreatedBy:        queue.CreatedBy.String(),
		})
	}

	return queues, nil
}

func (u QueueUsercase) GetQueueByID(ctx *gin.Context, id string) (*model.GetQueuesResponse, error) {
	queueId, err := uuid.Parse(id)
	if err != nil {
		return nil, err
	}

	queueRepo, err := u.repo.GetQueueByID(ctx, queueId)
	if err != nil {
		return nil, err
	}

	queue := &model.GetQueuesResponse{
		QueueID:          queueRepo.ID,
		QueueName:        queueRepo.Name,
		QueueDescription: queueRepo.Description,
		CreatedAt:        queueRepo.CreatedAt,
		CreatedBy:        queueRepo.CreatedBy.String(),
	}

	return queue, nil
}
