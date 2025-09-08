package updatecaseupdater

import (
	"case-management/internal/domain/repository"
	"context"
	"fmt"

	"github.com/google/uuid"
)

type ChangeCustomerInfoUpdater struct {
	repo repository.CaseRepository
}

func NewChangeCustomerInfoUpdater(repo repository.CaseRepository) *ChangeCustomerInfoUpdater {
	return &ChangeCustomerInfoUpdater{repo: repo}
}

func (u *ChangeCustomerInfoUpdater) Update(ctx context.Context, caseID uuid.UUID, data map[string]interface{}) error {
	if data["currentInfo"] == nil {
		return fmt.Errorf("currentInfo is required")
	}

	if data["newInfo"] == nil {
		return fmt.Errorf("newInfo is required")
	}

	return nil
}
