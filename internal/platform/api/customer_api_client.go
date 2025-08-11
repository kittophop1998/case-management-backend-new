package api

import (
	"bytes"
	"case-management/internal/domain/model"
	externalmodel "case-management/internal/domain/model/external_api"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/go-playground/validator/v10"
)

var (
	ErrRequestTimeout = errors.New("REQUEST_TIMEOUT")
	ErrBadRequest     = errors.New("BAD_REQUEST")
	ErrNotFound       = errors.New("NOT_FOUND")
	ErrInternal       = errors.New("INTERNAL_SERVER_ERROR")
)

type dashboardAPIClient struct {
	baseURL    string
	httpClient *http.Client
}

func NewDashboardAPIClient(baseURL string) *dashboardAPIClient {
	return &dashboardAPIClient{
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

func (c *dashboardAPIClient) FindByAeonID(ctx context.Context, aeonID string) (*model.GetCustomerInfoResponse, error) {

	url := fmt.Sprintf("%s/GetCustomerInfo", c.baseURL)

	payload := externalmodel.ConnectorGetCustomerInfoRequest{
		UserRef: aeonID,
		Mode:    "F",
	}

	reqBody, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("cannot marshal request body: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, io.NopCloser(bytes.NewReader(reqBody)))
	if err != nil {
		return nil, fmt.Errorf("cannot create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			return nil, ErrRequestTimeout
		}
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("cannot read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		var apiErr externalmodel.ConnectorErrorResponse
		if json.Unmarshal(bodyBytes, &apiErr) == nil && apiErr.ErrorCode == "COM001" {
			return nil, ErrBadRequest
		}
		return nil, fmt.Errorf("API returned status: %d", resp.StatusCode)
	}

	var apiResp externalmodel.ConnectorGetCustomerInfoResponse
	if err := json.Unmarshal(bodyBytes, &apiResp); err != nil {
		return nil, fmt.Errorf("cannot parse response: %w", err)
	}

	if err := validator.New().Struct(apiResp); err != nil {
		return nil, ErrNotFound
	}

	customer := &model.GetCustomerInfoResponse{
		CustomerNameEng: apiResp.CustomerNameEng,
		CustomerNameTH:  apiResp.CustomerNameTH,
		NationalID:      apiResp.CustID,
		MobileNO:        trimMobilePrefix(apiResp.MobileNO),
		MailTo:          apiResp.MailTo,
		MailToAddress:   apiResp.MailToAddress,
	}

	if customer.MailTo == "" {
		if apiResp.MailTo == "1" {
			customer.MailTo = "Home"
			customer.MailToAddress = fmt.Sprintf("%s %d", apiResp.HomeAddress, apiResp.HomeZip)
		} else {
			customer.MailTo = "Office"
			customer.MailToAddress = fmt.Sprintf("%s %s %d", apiResp.OfficeName, apiResp.OfficeAddress, apiResp.OfficeZip)
		}
	}

	return customer, nil
}

func trimMobilePrefix(mobile string) string {
	if len(mobile) > 1 {
		return mobile[1:]
	}
	return mobile
}
