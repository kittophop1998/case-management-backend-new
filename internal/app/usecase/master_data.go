package usecase

import (
	"case-management/internal/domain/repository"

	"github.com/gin-gonic/gin"
)

type MasterDataUseCase struct {
	repo repository.MasterDataRepository
}

func NewMasterDataUseCase(repo repository.MasterDataRepository) *MasterDataUseCase {
	return &MasterDataUseCase{repo: repo}
}

func (m *MasterDataUseCase) GetAllLookups(ctx *gin.Context) (map[string]interface{}, error) {
	data, err := m.repo.FindAll(ctx)
	if err != nil {
		return nil, err
	}
	return data, nil
}
