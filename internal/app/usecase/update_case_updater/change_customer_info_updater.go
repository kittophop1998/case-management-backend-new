package updatecaseupdater

import (
	"case-management/internal/domain/model"
	"case-management/internal/domain/repository"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type ChangeCustomerInfoUpdater struct {
	repo repository.CaseRepository
}

func NewChangeCustomerInfoUpdater(repo repository.CaseRepository) *ChangeCustomerInfoUpdater {
	return &ChangeCustomerInfoUpdater{repo: repo}
}

func (u *ChangeCustomerInfoUpdater) Update(ctx *gin.Context, caseID uuid.UUID, data map[string]interface{}) error {
	changeCustomerInfo := model.ChangeCustomerInfo{
		CurrentInfo: data["currentInfo"].(string),
		NewInfo:     data["newInfo"].(string),
	}
	fmt.Println("ChangeCustomerInfoUpdater called with:", changeCustomerInfo)

	return nil
}
