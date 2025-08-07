package usecase

import (
	"case-management/internal/domain/model"
	"case-management/internal/domain/repository"

	"github.com/gin-gonic/gin"
)

type UserUseCase interface {
	FindAll(ctx *gin.Context) ([]model.User, error)
}

type userUseCase struct {
	repo repository.UserRepository
}

func NewUserUsecase(repo repository.UserRepository) UserUseCase {
	return &userUseCase{repo: repo}
}

func (u *userUseCase) FindAll(ctx *gin.Context) ([]model.User, error) {
	return u.repo.FindAll(ctx)
}
