package model

import (
	"time"

	"github.com/google/uuid"
)

type Queues struct {
	ID          uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey" json:"id"`
	Name        string    `gorm:"type:varchar(100);not null" json:"name"`
	Description string    `gorm:"type:text" json:"description"`
	CreatedBy   uuid.UUID `gorm:"type:uuid" json:"createdBy"`
	CreatedAt   time.Time `gorm:"type:timestamptz;" json:"createdAt"`
	UpdatedBy   uuid.UUID `gorm:"type:uuid" json:"updatedBy"`
	UpdatedAt   time.Time `gorm:"type:timestamptz;" json:"updatedAt"`
	DeletedAt   time.Time `gorm:"type:timestamptz" json:"deletedAt"`
	DeletedBy   uuid.UUID `gorm:"type:uuid" json:"deletedBy"`
}

type QueueUsers struct {
	QueueID   uuid.UUID `gorm:"type:uuid;primaryKey" json:"queueId"`
	UserID    uuid.UUID `gorm:"type:uuid;primaryKey" json:"userId"`
	CreatedBy uuid.UUID `gorm:"type:uuid" json:"createdBy"`
	CreatedAt time.Time `gorm:"type:timestamptz;" json:"createdAt"`
	UpdatedBy uuid.UUID `gorm:"type:uuid" json:"updatedBy"`
	UpdatedAt time.Time `gorm:"type:timestamptz;" json:"updatedAt"`
	DeletedAt time.Time `gorm:"type:timestamptz" json:"deletedAt"`
	DeletedBy uuid.UUID `gorm:"type:uuid" json:"deletedBy"`
}

// ##### Response For Queue #####
type GetQueuesResponse struct {
	QueueID          string    `json:"queueId"`
	QueueName        string    `json:"queueName"`
	QueueDescription string    `json:"queueDescription"`
	CreatedAt        time.Time `json:"createdAt"`
	CreatedBy        string    `json:"createdBy"`
}

// ##### Request For Queue #####
type CreateQueueRequest struct {
	QueueName        string `json:"queueName" binding:"required"`
	QueueDescription string `json:"queueDescription"`
}

type UserManageInQueue struct {
	Users []string `json:"users" binding:"required"`
}
