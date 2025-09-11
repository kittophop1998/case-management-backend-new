package api

import (
	"bytes"
	"case-management/infrastructure/config"
	"case-management/internal/domain/model"
	externalmodel "case-management/internal/domain/model/external_api"
	"case-management/utils"
	"context"
	"crypto/tls"
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

type DashboardAPIClient struct {
	BaseURL    string
	TDURL      string
	httpClient *http.Client
	cfg        *config.Config
}

type DashboardAPIProxyClient struct {
	BaseURL    string
	TDURL      string
	httpClient *http.Client
	cfg        *config.Config
}

func NewDashboardAPIClient(cfg *config.Config) *DashboardAPIClient {

	defaultTr := http.DefaultTransport.(*http.Transport)
	transpot := defaultTr.Clone()
	transpot.TLSClientConfig = &tls.Config{InsecureSkipVerify: true}

	client := &http.Client{
		Timeout:   10 * time.Second,
		Transport: transpot,
	}

	return &DashboardAPIClient{
		BaseURL:    cfg.ConnectorAPIConfig.BaseURL,
		TDURL:      cfg.TDAPIConfig.BaseURL,
		httpClient: client,
		cfg:        cfg,
	}
}

func NewDashboardAPIProxyClient(cfg *config.Config) *DashboardAPIProxyClient {

	defaultTr := http.DefaultTransport.(*http.Transport)
	transpot := defaultTr.Clone()
	transpot.TLSClientConfig = &tls.Config{InsecureSkipVerify: true}

	proxyURLStr := "http://10.255.100.10:8080"
	proxyURL, err := url.Parse(proxyURLStr)
	if err != nil {
		log.Printf("[NewDashboardAPIClient] Invalid proxy URL: %v", err)
	}

	transpot.Proxy = http.ProxyURL(proxyURL)

	client := &http.Client{
		Timeout:   10 * time.Second,
		Transport: transpot,
	}

	return &DashboardAPIProxyClient{
		BaseURL:    cfg.ConnectorAPIConfig.BaseURL,
		TDURL:      cfg.TDAPIConfig.BaseURL,
		httpClient: client,
		cfg:        cfg,
	}
}

func (c *DashboardAPIClient) GetCustInfoByAeonID(ctx context.Context, reqData model.ConnectorCustomerInfoRequest) (*model.GetCustInfoResponse, http.Header, error) {
	url := fmt.Sprintf("%s/Api/Common/GetCustomerInfo", c.BaseURL)

	payload := reqData

	reqBody, err := json.Marshal(payload)
	if err != nil {
		return nil, nil, fmt.Errorf("cannot marshal request body: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(reqBody))
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create request: %w", err)
	}

	log.Println("Payload : ", payload)

	req.Header.Set("Content-Type", "application/json")

	// ใช้ SetHeadersFormContext เพื่อเซ็ต header จาก gin.Context
	utils.SetHeadersFormContext(ctx, req, []utils.CtxKey{
		utils.CtxKeyApisKey,
		utils.CtxKeyApiLang,
		utils.CtxKeyChannel,
		utils.CtxKeyDeviceOS,
		utils.CtxKeyRequestID,
	})

	log.Println("[API Client] => Request Headers:")
	for k, v := range req.Header {
		log.Printf("  %s: %v\n", k, v)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			return nil, nil, ErrRequestTimeout
		}
		return nil, nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, nil, fmt.Errorf("cannot read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		var apiErr externalmodel.ConnectorErrorResponse
		if json.Unmarshal(bodyBytes, &apiErr) == nil && apiErr.ErrorCode == "COM001" {
			return nil, nil, ErrBadRequest
		}
		return nil, nil, fmt.Errorf("API returned status: %d", resp.StatusCode)
	}

	var apiResp externalmodel.ConnectorGetCustomerInfoResponse
	if err := json.Unmarshal(bodyBytes, &apiResp); err != nil {
		return nil, nil, fmt.Errorf("cannot parse response: %w", err)
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

	return customer, resp.Header, nil
}

// GetCustProfileByAeonID retrieves customer profile information by Aeon ID.
func (c *DashboardAPIProxyClient) GetCustProfileByAeonID(ctx context.Context, aeonID string) (*model.GetCustProfileResponse, http.Header, error) {
	version := "2"
	token := "f3eedb05-f578-40c6-94c3-5a5ee60e9376,802e5979-3abf-41f5-b2b6-015af1146f03,d7743e18-320e-4647-bfcd-da0676c300e7,ef158183-1c6e-4559-82a7-4b6d629e0cf3,1861457b-8d6c-4797-9bdd-184d8ce4d01a"

	url, err := c.buildURL(version, token, aeonID)
	if err != nil {
		log.Println("[GetCustProfileByAeonID] Failed to build URL:", err)
		return nil, nil, err
	}

	log.Println("[GetCustProfileByAeonID] Request URL:", url)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		log.Println("[GetCustProfileByAeonID] Failed to create request:", err)
		if errors.Is(err, context.DeadlineExceeded) || errors.Is(err, context.Canceled) {
			return nil, nil, ErrRequestTimeout
		}
		return nil, nil, fmt.Errorf("cannot create request: %w", err)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		log.Println("[GetCustProfileByAeonID] Request failed:", err)
		if errors.Is(err, context.DeadlineExceeded) {
			return nil, nil, ErrRequestTimeout
		}
		return nil, nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		// อ่านแค่ขนาดจำกัดของ body ในกรณี error เพื่อป้องกันใช้หน่วยความจำมากเกินไป
		limitedBody, _ := io.ReadAll(io.LimitReader(resp.Body, 1024*10)) // 10 KB
		log.Printf("[GetCustProfileByAeonID] API returned status: %d, body: %s\n", resp.StatusCode, string(limitedBody))
		return nil, nil, fmt.Errorf("API returned status: %d, response body: %s", resp.StatusCode, string(limitedBody))
	}

	// ใช้ json.Decoder แทนการอ่าน body ทั้งหมด
	decoder := json.NewDecoder(resp.Body)

	var items []struct {
		Attributes map[string]string `json:"attributes"`
	}

	if err := decoder.Decode(&items); err != nil {
		log.Println("[GetCustProfileByAeonID] Failed to decode JSON:", err)
		return nil, nil, fmt.Errorf("cannot decode API response: %w", err)
	}

	if len(items) == 0 {
		log.Println("[GetCustProfileByAeonID] No items found in response")
		return nil, nil, ErrNotFound
	}

	respData := &model.GetCustProfileResponse{}

	for _, item := range items {
		attrs := item.Attributes

		// เอาข้อมูลจาก attributes มาแมปใส่ struct ทีละฟิลด์
		if v, ok := attrs["last_card_apply_date"]; ok {
			respData.LastCardApplyDate = v
		}
		if v, ok := attrs["customer_sentiment"]; ok {
			respData.CustomerSentiment = v
		}
		if v, ok := attrs["phone_no_last_update_date"]; ok {
			respData.PhoneNoLastUpdateDate = v
		}
		if v, ok := attrs["last_increase_credit_limit_update"]; ok {
			respData.LastIncreaseCreditLimitUpdate = v
		}
		if v, ok := attrs["last_reduce_credit_limit_update"]; ok {
			respData.LastReduceCreditLimitUpdate = v
		}
		if v, ok := attrs["last_income_update"]; ok {
			respData.LastIncomeUpdate = v
		}
		if v, ok := attrs["suggested_action"]; ok {
			respData.SuggestedAction = v
		}
		if v, ok := attrs["type_of_job"]; ok {
			respData.TypeOfJob = v
		}
		if v, ok := attrs["marital_status"]; ok {
			respData.MaritalStatus = v
		}
		if v, ok := attrs["gender"]; ok {
			respData.Gender = v
		}
		if v, ok := attrs["last_e_statement_sent_date"]; ok {
			respData.LastEStatementSentDate = v
		}
		if v, ok := attrs["e_statement_sent_status"]; ok {
			respData.EStatementSentStatus = v
		}
		if v, ok := attrs["statement_channel"]; ok {
			respData.StatementChannel = v
		}
		if v, ok := attrs["consent_for_disclose"]; ok {
			respData.ConsentForDisclose = v
		}
		if v, ok := attrs["block_media"]; ok {
			respData.BlockMedia = v
		}
		if v, ok := attrs["consent_for_collect_use"]; ok {
			respData.ConsentForCollectUse = v
		}

		// อื่น ๆ เช่น payment_status, day_past_due, last_overdue_date
		if v, ok := attrs["payment_status"]; ok {
			respData.PaymentStatus = v
		}
		if v, ok := attrs["day_past_due"]; ok {
			respData.DayPastDue = v
		}
		if v, ok := attrs["last_overdue_date"]; ok {
			respData.LastOverdueDate = v
		}
	}

	log.Println("[GetCustProfileByAeonID] Successfully built profile response")

	return respData, resp.Header, nil
}

// GetCustSegmentByAeonID retrieves customer segment information by Aeon ID.
func (c *DashboardAPIProxyClient) GetCustSegmentByAeonID(ctx context.Context, aeonID string) (*model.GetCustSegmentResponse, http.Header, error) {

	version := "2"
	token := "0fc5ded4-4b84-40f3-83e5-9711310a87df,9b6245ff-a10e-4e49-b039-e43039c5c1d3"

	url, err := c.buildURL(version, token, aeonID)
	if err != nil {
		return nil, nil, err
	}

	log.Println("Path URL : ", url)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) || errors.Is(err, context.Canceled) {
			return nil, nil, ErrRequestTimeout
		}
		return nil, nil, fmt.Errorf("cannot create request: %w", err)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			return nil, nil, ErrRequestTimeout
		}
		return nil, nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, nil, fmt.Errorf("cannot read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, nil, fmt.Errorf("API returned status: %d, response body: %s", resp.StatusCode, string(bodyBytes))
	}

	var items []externalmodel.Item
	if err := json.Unmarshal(bodyBytes, &items); err != nil {
		return nil, nil, fmt.Errorf("cannot unmarshal API response into items array: %w", err)
	}

	if len(items) < 2 || len(items[0].Attributes) == 0 || len(items[1].Attributes) == 0 {
		return nil, nil, ErrNotFound
	}

	var attr1 externalmodel.Attributes
	if err := json.Unmarshal(items[0].Attributes, &attr1); err != nil {
		return nil, nil, fmt.Errorf("failed to unmarshal first attributes: %w", err)
	}

	var attr2 externalmodel.Attributes
	if err := json.Unmarshal(items[1].Attributes, &attr2); err != nil {
		return nil, nil, fmt.Errorf("failed to unmarshal second attributes: %w", err)
	}

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
	if strings.TrimSpace(attr1.VvipCustomerPosition) != "" {
		custsegment.ComplaintGroup = attr1.VvipCustomerGroupFlag + " - " + attr1.VvipCustomerPosition
	} else {
		custsegment.ComplaintGroup = attr1.VvipCustomerGroupFlag
	}

	return custsegment, resp.Header, nil
}

// GetCustSuggestionByAeonID retrieves customer suggestions by Aeon ID.
func (c *DashboardAPIProxyClient) GetCustSuggestionByAeonID(ctx context.Context, aeonID string) (*model.GetCustSuggestionResponse, http.Header, error) {

	version := "2"
	token := "51fb37b4-618e-4c66-8f15-9177e7f646d3,4071bc51-4dbb-449f-9d71-f929936fed1b"

	url, err := c.buildURL(version, token, aeonID)
	if err != nil {
		return nil, nil, err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) || errors.Is(err, context.Canceled) {
			return nil, nil, ErrRequestTimeout
		}
		return nil, nil, fmt.Errorf("cannot create request: %w", err)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			return nil, nil, ErrRequestTimeout
		}
		return nil, nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, nil, fmt.Errorf("cannot read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, nil, fmt.Errorf("API returned status: %d, response body: %s", resp.StatusCode, string(bodyBytes))
	}

	var items []externalmodel.Item
	if err := json.Unmarshal(bodyBytes, &items); err != nil {
		return nil, nil, fmt.Errorf("cannot unmarshal API response into items array: %w", err)
	}

	if len(items) < 2 || len(items[0].Attributes) == 0 || len(items[1].Attributes) == 0 {
		return nil, nil, ErrNotFound
	}

	var attr1 externalmodel.Attributes
	if err := json.Unmarshal(items[0].Attributes, &attr1); err != nil {
		return nil, nil, fmt.Errorf("failed to unmarshal first attributes: %w", err)
	}

	var attr2 externalmodel.Attributes
	if err := json.Unmarshal(items[1].Attributes, &attr2); err != nil {
		return nil, nil, fmt.Errorf("failed to unmarshal second attributes: %w", err)
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

	return suggestion, resp.Header, nil
}

// buildURL creates the complete API URL with all necessary query parameters
func (c *DashboardAPIProxyClient) buildURL(version, token, aeonID string) (string, error) {
	baseURL, err := url.Parse(c.TDURL)
	if err != nil {
		return "", fmt.Errorf("invalid base URL: %w", err)
	}

	query := fmt.Sprintf("version=%s&token=%s&key.aeon_id=%s", version, token, aeonID)

	fullURL := fmt.Sprintf("%s?%s", baseURL, query)
	return fullURL, nil

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
