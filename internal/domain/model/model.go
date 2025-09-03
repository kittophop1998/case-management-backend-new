package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// ##### Customer #####
type PermissionWithRolesResponse struct {
	Permission string   `json:"permission"`
	Name       string   `json:"name"`
	Roles      []string `json:"roles"`
}

type Model struct {
	ID        uuid.UUID      `gorm:"primaryKey;default:uuid_generate_v4()" json:"id"`
	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-" swaggerignore:"true"`
}
