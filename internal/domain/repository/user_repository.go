package repository

import (
	"case-management/internal/domain/model"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type UserRepository interface {
	GetAll(c *gin.Context) ([]*model.User, error)
	GetById(c *gin.Context, id uuid.UUID) (*model.User, error)
	GetByUsername(c *gin.Context, username string) (*model.User, error)
	Create(c *gin.Context, user *model.CreateUpdateUserRequest) (uuid.UUID, error)
	Update(c *gin.Context, userID uuid.UUID, input model.CreateUpdateUserRequest) error
	Count(c *gin.Context) (int, error)
	CountWithFilter(c *gin.Context, filter model.UserFilter) (int, error)
}
