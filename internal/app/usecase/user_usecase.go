package usecase

import (
	"case-management/internal/domain/model"
	"case-management/internal/domain/repository"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type UserUseCase struct {
	repo repository.UserRepository
}

func NewUserUseCase(repo repository.UserRepository) *UserUseCase {
	return &UserUseCase{repo: repo}
}

func (u *UserUseCase) GetAll(ctx *gin.Context) ([]*model.User, error) {
	return u.repo.GetAll(ctx)
}

func (u *UserUseCase) GetById(ctx *gin.Context, id uuid.UUID) (*model.User, error) {
	return u.repo.GetById(ctx, id)
}

func (u *UserUseCase) GetByUsername(ctx *gin.Context, username string) (*model.User, error) {
	return u.repo.GetByUsername(ctx, username)
}

func (u *UserUseCase) Create(ctx *gin.Context, user *model.CreateUpdateUserRequest) (uuid.UUID, error) {
	return u.repo.Create(ctx, user)
}

func (u *UserUseCase) Update(ctx *gin.Context, userID uuid.UUID, input model.CreateUpdateUserRequest) error {
	return u.repo.Update(ctx, userID, input)
}

func (u *UserUseCase) Count(ctx *gin.Context) (int, error) {
	return u.repo.Count(ctx)
}
