package repository

import (
	"case-management/internal/domain/model"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type UserRepository interface {
	GetAll(ctx *gin.Context) ([]*model.User, error)
	GetById(ctx *gin.Context, id uuid.UUID) (*model.User, error)
	GetByUsername(ctx *gin.Context, username string) (*model.User, error)
	Create(ctx *gin.Context, user *model.CreateUpdateUserRequest) (uuid.UUID, error)
	Update(ctx *gin.Context, id uuid.UUID, input model.CreateUpdateUserRequest) error
	Count(ctx *gin.Context) (int, error)
	CountWithFilter(ctx *gin.Context, filter model.UserFilter) (int, error)
}
