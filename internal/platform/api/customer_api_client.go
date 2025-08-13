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
	"net/url"
	"strings"
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

func (c *dashboardAPIClient) GetCustInfoByAeonID(ctx context.Context, aeonID string) (*model.GetCustInfoResponse, error) {

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

	customer := &model.GetCustInfoResponse{
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
func (c *dashboardAPIClient) GetCustProfileByAeonID(ctx context.Context, aeonID string) (*model.GetCustProfileResponse, error) {

	version := "2"
	token := "f3eedb05-f578-40c6-94c3-5a5ee60e9376,802e5979-3abf-41f5-b2b6-015af1146f03,d7743e18-320e-4647-bfcd-da0676c300e7,ef158183-1c6e-4559-82a7-4b6d629e0cf3,1861457b-8d6c-4797-9bdd-184d8ce4d01a"

	url, err := c.buildURL(version, token, aeonID)
	if err != nil {
		return nil, err
	}

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
	var custprofile model.GetCustProfileResponse

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

// GetCustSegmentByAeonID retrieves customer segment information by Aeon ID.
func (c *dashboardAPIClient) GetCustSegmentByAeonID(ctx context.Context, aeonID string) (*model.GetCustSegmentResponse, error) {

	version := "2"
	token := "0fc5ded4-4b84-40f3-83e5-9711310a87df,9b6245ff-a10e-4e49-b039-e43039c5c1d3"

	url, err := c.buildURL(version, token, aeonID)
	if err != nil {
		return nil, err
	}

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

	if len(items) < 2 || len(items[0].Attributes) == 0 || len(items[1].Attributes) == 0 {
		return nil, ErrNotFound
	}

	var attr1 externalmodel.Attributes
	if err := json.Unmarshal(items[0].Attributes, &attr1); err != nil {
		return nil, fmt.Errorf("failed to unmarshal first attributes: %w", err)
	}

	var attr2 externalmodel.Attributes
	if err := json.Unmarshal(items[1].Attributes, &attr2); err != nil {
		return nil, fmt.Errorf("failed to unmarshal second attributes: %w", err)
	}

	// Map the unmarshaled data to the final response struct with a more compact syntax.
	custsegment := &model.GetCustSegmentResponse{
		Sweetheart:      attr1.SweetheartCustomerGroupFlag,
		CustomerType:    attr2.CustomerType,
		MemberStatus:    attr2.MemberStatus,
		CustomerSegment: attr2.CbaSegment,
		UpdateData:      formatDateIfValid(attr1.UpdateData),
	}

	// ComplaintLevel
	if attr1.CustomerLevel != "NORMAL" && strings.TrimSpace(attr1.CustomerLevel) != "" {
		custsegment.ComplaintLevel = "Complaint Level: " + attr1.CustomerLevel
	} else {
		custsegment.ComplaintLevel = attr1.CustomerLevel
	}

	// CustomerGroup
	if strings.TrimSpace(attr1.VvipCustomerPosition) != "" {
		custsegment.CustomerGroup = attr1.VvipCustomerGroupFlag + " - " + attr1.VvipCustomerPosition
	} else {
		custsegment.CustomerGroup = attr1.VvipCustomerGroupFlag
	}

	// ComplaintGroup
	if strings.TrimSpace(attr2.ComplaintTopic) != "" {
		custsegment.ComplaintGroup = attr2.ComplaintGroup + " (" + attr2.ComplaintTopic + ")"
	} else {
		custsegment.ComplaintGroup = attr2.ComplaintGroup
	}

	return custsegment, nil
}

// GetCustSuggestionByAeonID retrieves customer suggestions by Aeon ID.
func (c *dashboardAPIClient) GetCustSuggestionByAeonID(ctx context.Context, aeonID string) (*model.GetCustSuggestionResponse, error) {

	version := "2"
	token := "51fb37b4-618e-4c66-8f15-9177e7f646d3,4071bc51-4dbb-449f-9d71-f929936fed1b"

	url, err := c.buildURL(version, token, aeonID)
	if err != nil {
		return nil, err
	}

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

	if len(items) < 2 || len(items[0].Attributes) == 0 || len(items[1].Attributes) == 0 {
		return nil, ErrNotFound
	}

	var attr1 externalmodel.Attributes
	if err := json.Unmarshal(items[0].Attributes, &attr1); err != nil {
		return nil, fmt.Errorf("failed to unmarshal first attributes: %w", err)
	}

	var attr2 externalmodel.Attributes
	if err := json.Unmarshal(items[1].Attributes, &attr2); err != nil {
		return nil, fmt.Errorf("failed to unmarshal second attributes: %w", err)
	}

	// Initialize the final response object.
	suggestion := &model.GetCustSuggestionResponse{}

	// Map suggested cards from attr1.
	nameOfCards := strings.Split(attr1.NameOfCards, ",")
	for _, card := range nameOfCards {
		suggestion.SuggestCards = append(suggestion.SuggestCards, strings.TrimSpace(card))
	}

	// Map suggested promotions from attr2.
	for _, row := range attr2.PromotionArray {
		var promotion model.GetCustSuggestionPromotionResponse
		if len(row) < 7 {
			log.Printf("Incomplete promotion data row: %v", row)
			continue
		}

		promotion.PromotionCode = row[0]
		promotion.PromotionName = row[1]
		promotion.PromotionDetails = row[2]
		promotion.Action = row[5]

		promotion.PromotionResultTimestamp = formatDateTimeIfValid(row[6])

		periodList := strings.Split(row[3], " ")
		for i, period := range periodList {
			periodList[i] = formatDateIfValid(period)
		}
		promotion.Period = strings.Join(periodList, " - ")

		eligibleCards := strings.Split(row[4], ",")
		for _, card := range eligibleCards {
			promotion.EligibleCard = append(promotion.EligibleCard, strings.TrimSpace(card))
		}

		suggestion.SuggestPromotions = append(suggestion.SuggestPromotions, promotion)
	}

	return suggestion, nil
}

// buildURL creates the complete API URL with all necessary query parameters
func (c *dashboardAPIClient) buildURL(version, token, aeonID string) (string, error) {
	u, err := url.Parse(c.baseURL)
	if err != nil {
		return "", fmt.Errorf("invalid base URL: %w", err)
	}

	queryParams := u.Query()
	queryParams.Set("version", version)
	queryParams.Set("token", token)
	queryParams.Set("key.aeon_id", aeonID)
	u.RawQuery = queryParams.Encode()

	return u.String(), nil
}

func trimMobilePrefix(mobile string) string {
	if len(mobile) > 1 {
		return mobile[1:]
	}
	return mobile
}

// formatDateIfValid is a helper function to format a date string if it's not empty.
func formatDateIfValid(dateStr string) string {
	if strings.TrimSpace(dateStr) == "" {
		return dateStr
	}
	dateTime, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		log.Printf("Error parsing date '%s': %v", dateStr, err)
		return dateStr
	}
	return dateTime.Format("02 Jan 2006")
}

// formatDateTimeIfValid is a helper function to format a datetime string if it's not empty.
func formatDateTimeIfValid(dateTimeStr string) string {
	if strings.TrimSpace(dateTimeStr) == "" {
		return dateTimeStr
	}
	dateTime, err := time.Parse("2006-01-02 15.04.05", dateTimeStr)
	if err != nil {
		log.Printf("Error parsing datetime '%s': %v", dateTimeStr, err)
		return dateTimeStr
	}
	return dateTime.Format("02 Jan 2006, 15.04")
}
