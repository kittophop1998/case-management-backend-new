package usecase

import (
	"case-management/internal/domain/model"
	"case-management/internal/domain/repository"
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type QueueUsecase struct {
	auditLogger repository.AuditLogRepository
	repo        repository.QueueRepository
}

func NewQueueUsecase(
	auditLogger repository.AuditLogRepository,
	repo repository.QueueRepository,
) *QueueUsecase {
	return &QueueUsecase{auditLogger: auditLogger, repo: repo}
}

func (u QueueUsecase) GetQueues(ctx *gin.Context, page, limit int, queueName string) ([]*model.GetQueuesResponse, int, error) {
	var queues []*model.GetQueuesResponse
	offset := (page - 1) * limit
	queuesRepo, total, err := u.repo.GetQueues(ctx, offset, limit, queueName)
	if err != nil {
		return nil, 0, err
	}

	if len(queuesRepo) == 0 {
		queues = []*model.GetQueuesResponse{}
		return queues, 0, nil
	}

	for _, queue := range queuesRepo {
		queues = append(queues, &model.GetQueuesResponse{
			QueueID:          queue.ID.String(),
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
		QueueID:          queueRepo.ID.String(),
		QueueName:        queueRepo.Name,
		QueueDescription: queueRepo.Description,
		CreatedAt:        queueRepo.CreatedAt,
		CreatedBy:        queueRepo.CreatedBy.String(),
	}

	return queue, nil
}

func (u QueueUsecase) CreateQueue(ctx *gin.Context, createdBy string, input *model.CreateQueueRequest) (uuid.UUID, error) {
	// ##### Check if queue already exists #####
	isExisting := u.repo.IsExistingQueue(ctx, input.QueueName)
	if isExisting {
		return uuid.Nil, fmt.Errorf("queue with name %q already exists", input.QueueName)
	}

	// #####  Create Queue #####
	var createdByID uuid.UUID
	if id, err := uuid.Parse(createdBy); err == nil {
		createdByID = id
	}

	queue := &model.Queues{
		Name:        input.QueueName,
		Description: input.QueueDescription,
		CreatedAt:   time.Now(),
		CreatedBy:   createdByID,
		UpdatedBy:   createdByID,
	}

	queueID, err := u.repo.CreateQueue(ctx, queue)
	if err != nil {
		return uuid.Nil, err
	}

	// #####  CreateQueueUser #####
	var queueUsers []*model.QueueUsers
	for _, user := range input.QueueUsers {
		var userID uuid.UUID
		if id, err := uuid.Parse(user); err == nil {
			userID = id
		}

		queueUsers = append(queueUsers, &model.QueueUsers{
			QueueID:   queueID,
			UserID:    userID,
			CreatedAt: time.Now(),
			CreatedBy: createdByID,
			UpdatedBy: createdByID,
		})
	}

	if queueUsers != nil {
		if err := u.repo.CreateQueueUser(ctx, queueUsers); err != nil {
			return uuid.Nil, err
		}
	}

	u.auditLogger.LogAction(ctx, model.AuditLogs{
		Action:    "create_queue",
		Module:    "queue",
		UserID:    createdByID,
		CreatedAt: time.Now(),
	})

	return queueID, nil
}

func (u QueueUsecase) UpdateQueue(ctx *gin.Context, updatedBy string, input *model.UpdateQueueRequest) error {
	var queueID uuid.UUID
	if id, err := uuid.Parse(input.QueueID); err != nil {
		queueID = id
	}

	queue, err := u.repo.GetQueueByID(ctx, queueID)
	if err != nil {
		return fmt.Errorf("failed to get queue by ID: %w", err)
	}

	if queue == nil {
		return fmt.Errorf("queue with ID %q does not exist", input.QueueID)
	}

	queueToSave := &model.Queues{
		ID:          queueID,
		Name:        input.QueueName,
		Description: input.QueueDescription,
		UpdatedAt:   time.Now(),
		UpdatedBy:   queue.UpdatedBy,
	}

	if err := u.repo.UpdateQueue(ctx, queueToSave); err != nil {
		return fmt.Errorf("failed to update queue: %w", err)
	}

	var queueUsersToSave []*model.QueueUsers
	for _, user := range input.UsersAdd {
		var userID uuid.UUID
		if id, err := uuid.Parse(user); err == nil {
			userID = id
		}

		queueUsersToSave = append(queueUsersToSave, &model.QueueUsers{
			QueueID:   queueID,
			UserID:    userID,
			CreatedAt: time.Now(),
			CreatedBy: queue.UpdatedBy,
			UpdatedBy: queue.UpdatedBy,
		})
	}

	if err := u.repo.UpdateQueueUser(ctx, queueID, queueUsersToSave, input.UsersDel); err != nil {
		return fmt.Errorf("failed to update queue users: %w", err)
	}

	return nil
}
