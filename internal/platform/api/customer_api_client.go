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
	"log"
	"net/http"
	"time"
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

func (c *dashboardAPIClient) GetCustInfoByAeonID(ctx context.Context, aeonID string) (*model.GetCustomerInfoResponse, error) {

	url := fmt.Sprintf("%s/Common/GetCustomerInfo", c.baseURL)

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

// GetCustProfileByAeonID retrieves customer profile information by Aeon ID.
func (c *dashboardAPIClient) GetCustProfileByAeonID(ctx context.Context, aeonID string) (*model.GetCustomerProfileResponse, error) {

	url := fmt.Sprintf("%s/customerprofile?aeonID=%s", c.baseURL, aeonID)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) || errors.Is(err, context.Canceled) {
			return nil, ErrRequestTimeout
		}
		return nil, fmt.Errorf("cannot create request: %w", err)
	}

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
		return nil, fmt.Errorf("API returned status: %d, response body: %s", resp.StatusCode, string(bodyBytes))
	}

	var items []externalmodel.Item
	if err := json.Unmarshal(bodyBytes, &items); err != nil {
		return nil, fmt.Errorf("cannot unmarshal API response into items array: %w", err)
	}

	if len(items) == 0 {
		return nil, ErrNotFound
	}

	// Initialize a new customer profile response object
	var custprofile model.GetCustomerProfileResponse

	profilesMap := make(map[string]externalmodel.Profile)

	for _, item := range items {
		if len(item.Attributes) == 0 || string(item.Attributes) == "{}" {
			continue
		}

		// Unmarshal based on the profile type
		switch item.Key.AeonID {
		case "profile1":
			var profile externalmodel.Profile1Response
			if err := json.Unmarshal(item.Attributes, &profile); err != nil {
				return nil, fmt.Errorf("failed to unmarshal Profile1: %w", err)
			}
			profilesMap["profile1"] = profile
		case "profile2":
			var profile externalmodel.Profile2Response
			if err := json.Unmarshal(item.Attributes, &profile); err != nil {
				return nil, fmt.Errorf("failed to unmarshal Profile2: %w", err)
			}
			profilesMap["profile2"] = profile
		case "profile3":
			var profile externalmodel.Profile3Response
			if err := json.Unmarshal(item.Attributes, &profile); err != nil {
				return nil, fmt.Errorf("failed to unmarshal Profile3: %w", err)
			}
			profilesMap["profile3"] = profile
		case "profile4":
			var profile externalmodel.Profile4Response
			if err := json.Unmarshal(item.Attributes, &profile); err != nil {
				return nil, fmt.Errorf("failed to unmarshal Profile4: %w", err)
			}
			profilesMap["profile4"] = profile
		case "profile5":
			var profile externalmodel.Profile5Response
			if err := json.Unmarshal(item.Attributes, &profile); err != nil {
				return nil, fmt.Errorf("failed to unmarshal Profile5: %w", err)
			}
			profilesMap["profile5"] = profile
		}
	}

	// Assign the unmarshaled data to the final response struct
	if p1, ok := profilesMap["profile1"].(externalmodel.Profile1Response); ok {
		custprofile.LastCardApplyDate = formatDateIfValid(p1.LastCardApplyDate)
		custprofile.PhoneNoLastUpdateDate = formatDateIfValid(p1.PhoneNoLastUpdateDate)
	}

	if p2, ok := profilesMap["profile2"].(externalmodel.Profile2Response); ok {
		custprofile.Gender = p2.Gender
		custprofile.MaritalStatus = p2.MaritalStatus
		custprofile.TypeOfJob = p2.TypeOfJob
	}

	if p3, ok := profilesMap["profile3"].(externalmodel.Profile3Response); ok {
		custprofile.LastIncreaseCreditLimitUpdate = formatDateIfValid(p3.LastIncreaseCreditLimitUpdate)
		custprofile.LastIncomeUpdate = formatDateIfValid(p3.LastIncomeUpdate)
		custprofile.LastReduceCreditLimitUpdate = formatDateIfValid(p3.LastReduceCreditLimitUpdate)
		custprofile.SuggestedAction = p3.SuggestedAction
	}

	if p4, ok := profilesMap["profile4"].(externalmodel.Profile4Response); ok {
		custprofile.LastEStatementSentDate = formatDateIfValid(p4.LastEStatementSentDate)
		custprofile.EStatementSentStatus = p4.EStatementSentStatus
		custprofile.StatementChannel = p4.StatementChannel
	}

	if p5, ok := profilesMap["profile5"].(externalmodel.Profile5Response); ok {
		custprofile.ConsentForCollectUse = p5.ConsentForCollectUse
		custprofile.BlockMedia = p5.BlockMedia
		custprofile.ConsentForDisclose = p5.ConsentForDisclose
	}

	return &custprofile, nil
}

func trimMobilePrefix(mobile string) string {
	if len(mobile) > 1 {
		return mobile[1:]
	}
	return mobile
}

func formatDateIfValid(dateString string) string {
	if dateString == "" || dateString == "No update" {
		return dateString
	}
	dateTime, err := time.Parse("2006-01-02", dateString)
	if err != nil {
		log.Printf("error parsing date string '%s': %v", dateString, err)
		return dateString
	}
	return dateTime.Format("02 Jan 2006")
}
