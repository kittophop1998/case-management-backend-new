package repository

import (
	"case-management/internal/domain/model"
	"context"
	"net/http"
)

type ConnectorRepository interface {
	GetCustInfoByAeonID(ctx context.Context, req model.ConnectorCustomerInfoRequest) (*model.GetCustInfoResponse, http.Header, error)
}

type TdRepository interface {
	GetCustProfileByAeonID(ctx context.Context, aeonID string) (*model.GetCustProfileResponse, http.Header, error)
	GetCustSegmentByAeonID(ctx context.Context, aeonID string) (*model.GetCustSegmentResponse, http.Header, error)
	GetCustSuggestionByAeonID(ctx context.Context, aeonID string) (*model.GetCustSuggestionResponse, http.Header, error)
}
