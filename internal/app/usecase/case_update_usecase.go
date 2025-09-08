package usecase

import (
	updatecaseupdater "case-management/internal/app/usecase/update_case_updater"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type UpdateCaseUseCase struct {
	Updaters map[string]updatecaseupdater.CaseUpdater
}

func NewUpdateCaseUseCase(updaters map[string]updatecaseupdater.CaseUpdater) *UpdateCaseUseCase {
	return &UpdateCaseUseCase{Updaters: updaters}
}

func (uc UpdateCaseUseCase) Execute(ctx *gin.Context, caseType string, caseID uuid.UUID, data map[string]interface{}) error {
	updater, ok := uc.Updaters[caseType]
	if !ok {
		return fmt.Errorf("unsupported case type: %s", caseType)
	}

	return updater.Update(ctx, caseID, data)
}
