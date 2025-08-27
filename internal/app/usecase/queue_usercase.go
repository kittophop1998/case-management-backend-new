package usecase

import (
	"case-management/internal/domain/model"
	"case-management/internal/domain/repository"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type QueueUsecase struct {
	repo repository.QueueRepository
}

func NewQueueUsecase(repo repository.QueueRepository) *QueueUsecase {
	return &QueueUsecase{repo: repo}
}

func (u QueueUsecase) GetQueues(ctx *gin.Context, page, limit int, queueName string) ([]*model.GetQueuesResponse, int, error) {
	var queues []*model.GetQueuesResponse
	queuesRepo, total, err := u.repo.GetQueues(ctx, page, limit, queueName)
	if err != nil {
		return nil, 0, err
	}

	if len(queuesRepo) == 0 {
		queues = []*model.GetQueuesResponse{}
		return queues, 0, nil
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

	return queues, total, nil
}

func (u QueueUsecase) GetQueueByID(ctx *gin.Context, id string) (*model.GetQueuesResponse, error) {
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
