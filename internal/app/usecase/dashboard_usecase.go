package usecase

import (
	"case-management/infrastructure/lib"
	"case-management/internal/domain/model"
	"case-management/internal/domain/repository"

	"github.com/gin-gonic/gin"
)

type DashboardUseCase struct {
	repo repository.DashboardRepository
}

func NewDashboardUseCase(repo repository.DashboardRepository) *DashboardUseCase {
	return &DashboardUseCase{repo: repo}
}

func (uc *DashboardUseCase) CustInfo(ctx *gin.Context, aeonID string) (*model.GetCustInfoResponse, error) {
	customer, err := uc.repo.GetCustInfoByAeonID(ctx, aeonID)
	if err != nil {
		details := map[string]string{
			"connector_api": "Connection issue from System-i",
		}
		lib.HandleError(ctx, lib.GatewayTimeout.WithDetails(details))
	}
	return customer, nil
}

func (uc *DashboardUseCase) CustProfile(ctx *gin.Context, aeonID string) (*model.GetCustProfileResponse, error) {
	custprofile, err := uc.repo.GetCustProfileByAeonID(ctx, aeonID)
	if err != nil {
		details := map[string]string{
			"cdp_api": "Connection issue from TD",
		}
		lib.HandleError(ctx, lib.GatewayTimeout.WithDetails(details))
	}
	return custprofile, nil
}

func (uc *DashboardUseCase) CustSegment(ctx *gin.Context, aeonID string) (*model.GetCustSegmentResponse, error) {
	custsegment, err := uc.repo.GetCustSegmentByAeonID(ctx, aeonID)
	if err != nil {
		details := map[string]string{
			"cdp_api": "Connection issue from TD",
		}
		lib.HandleError(ctx, lib.GatewayTimeout.WithDetails(details))
	}
	return custsegment, nil
}

func (uc *DashboardUseCase) CustSuggestion(ctx *gin.Context, aeonID string) (*model.GetCustSuggestionResponse, error) {
	custsuggestion, err := uc.repo.GetCustSuggestionByAeonID(ctx, aeonID)
	if err != nil {
		details := map[string]string{
			"cdp_api": "Connection issue from TD",
		}
		lib.HandleError(ctx, lib.GatewayTimeout.WithDetails(details))
	}
	return custsuggestion, nil
}
