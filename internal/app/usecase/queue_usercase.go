package usecase

import (
	"case-management/internal/domain/model"
	"case-management/internal/domain/repository"
	"case-management/utils"
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
	offset := (page - 1) * limit

	queuesRepo, total, err := u.repo.GetQueues(ctx, offset, limit, queueName)
	if err != nil {
		return nil, 0, err
	}

	var queues []*model.GetQueuesResponse
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
			CreatedBy:        utils.UserNameCenter(queue.CreatedUser),
		})
	}

	return queues, total, nil
}

func (u QueueUsecase) GetQueueByID(ctx *gin.Context, queueID uuid.UUID) (*model.GetQueuesResponse, error) {
	queueRepo, err := u.repo.GetQueueByID(ctx, queueID)
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

func (u QueueUsecase) CreateQueue(ctx *gin.Context, createdByID uuid.UUID, input *model.CreateQueueRequest) (uuid.UUID, error) {
	// ##### Check if queue already exists #####
	isExisting := u.repo.IsExistingQueue(ctx, input.QueueName)

	fmt.Println("Checking if queue exists:", isExisting)
	if isExisting {
		return uuid.Nil, fmt.Errorf("queue with name %q already exists", input.QueueName)
	}

	// #####  Create Queue #####
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

	u.auditLogger.LogAction(ctx, model.AuditLogs{
		Action:    "create_queue",
		Module:    "queue",
		UserID:    createdByID,
		CreatedAt: time.Now(),
	})

	return queueID, nil
}

func (u QueueUsecase) AddUserInQueue(ctx *gin.Context, createdByID uuid.UUID, queueID uuid.UUID, input model.UserManageInQueue) error {
	var queueUsers []*model.QueueUsers
	for _, user := range input.Users {
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
		if err := u.repo.AddQueueUser(ctx, queueUsers); err != nil {
			return err
		}
	}

	return nil
}

func (u QueueUsecase) UpdateQueueByID(ctx *gin.Context, updatedByID uuid.UUID, queueID uuid.UUID, input *model.UpdateQueueRequest) error {
	queue, err := u.repo.GetQueueByID(ctx, queueID)
	if err != nil {
		return fmt.Errorf("failed to get queue by ID: %w", err)
	}

	if queue == nil {
		return fmt.Errorf("queue with ID %q does not exist", queueID)
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

	return nil
}

func (u QueueUsecase) DeleteUsersInQueue(ctx *gin.Context, deletedByID uuid.UUID, queueID uuid.UUID, input model.UserManageInQueue) error {
	queue, err := u.repo.GetQueueByID(ctx, queueID)
	if err != nil {
		return fmt.Errorf("failed to get queue by ID: %w", err)
	}

	if queue == nil {
		return fmt.Errorf("queue with ID %q does not exist", queueID)
	}

	var userIDs []uuid.UUID
	for _, user := range input.Users {
		if id, err := uuid.Parse(user); err == nil {
			userIDs = append(userIDs, id)
		}
	}

	return u.repo.DeleteQueueUser(ctx, queueID, userIDs)
}
