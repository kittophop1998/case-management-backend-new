package usecase

import (
	"case-management/internal/domain/model"
	"case-management/internal/domain/repository"
	"case-management/pkg/api_errors"

	"github.com/gin-gonic/gin"
)

type DashboardUseCase struct {
	repo repository.DashboardRepository
}

func NewDashboardUseCase(repo repository.DashboardRepository) *DashboardUseCase {
	return &DashboardUseCase{repo: repo}
}

func (uc *DashboardUseCase) CustomerInfo(ctx *gin.Context, aeonID string) (*model.GetCustomerInfoResponse, error) {
	customer, err := uc.repo.GetCustInfoByAeonID(ctx, aeonID)
	if err != nil {
		details := map[string]string{
			"connector_api": "Connection issue from System-i",
		}
		appErr := api_errors.NewAppError(
			api_errors.ErrGatewayTimeout.Code,
			api_errors.ErrGatewayTimeout.Message,
			api_errors.ErrGatewayTimeout.HTTPStatus,
			details,
		)
		return nil, appErr
	}
	return customer, nil
}

func (uc *DashboardUseCase) CustomerProfile(ctx *gin.Context, aeonID string) (*model.GetCustomerProfileResponse, error) {
	custprofile, err := uc.repo.GetCustProfileByAeonID(ctx, aeonID)
	if err != nil {
		details := map[string]string{
			"connector_api": "Connection issue from System-i",
		}
		appErr := api_errors.NewAppError(
			api_errors.ErrGatewayTimeout.Code,
			api_errors.ErrGatewayTimeout.Message,
			api_errors.ErrGatewayTimeout.HTTPStatus,
			details,
		)
		return nil, appErr
	}
	return custprofile, nil
}
