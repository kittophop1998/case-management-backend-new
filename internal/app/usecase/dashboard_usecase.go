package usecase

import (
	"case-management/infrastructure/lib"
	"case-management/internal/domain/model"
	"case-management/internal/domain/repository"
	"case-management/internal/platform/api"
	"case-management/utils"
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type DashboardUseCase struct {
	repo            repository.DashboardRepository
	tdRepo          repository.TdRepository
	dashboardClient *api.DashboardAPIClient
	tdClient        *api.DashboardAPIProxyClient
	connectorRepo   repository.ConnectorRepository
}

func NewDashboardUseCase(repo repository.DashboardRepository, dashboardClient *api.DashboardAPIClient, tdClient *api.DashboardAPIProxyClient, connectorRepo repository.ConnectorRepository, tdRepo repository.TdRepository) *DashboardUseCase {
	return &DashboardUseCase{repo: repo, dashboardClient: dashboardClient, tdClient: tdClient, connectorRepo: connectorRepo, tdRepo: tdRepo}
}

// func (uc *DashboardUseCase) CustInfo(ctx context.Context, ginCtx *gin.Context, req model.ConnectorCustomerInfoRequest) (*model.GetCustInfoResponse, error) {
// 	log.Printf("[UseCase] => Called with request: %+v\n", req)

// 	// ===== Extract Request Header for logging =====
// 	var reqHeader string
// 	if ginCtx != nil && ginCtx.Request != nil {
// 		if headerBytes, err := json.Marshal(utils.FlattenHeader(ginCtx.Request.Header)); err == nil {
// 			reqHeader = string(headerBytes)
// 		}
// 	}

// 	// ===== Get Request ID from Context or Header =====
// 	requestID, _ := ctx.Value(utils.CtxKeyRequestID).(string)
// 	if requestID == "" && ginCtx != nil {
// 		requestID = ginCtx.GetHeader("X-Request-ID")
// 	}
// 	if requestID == "" {
// 		requestID = uuid.New().String()
// 	}

// 	// ===== Call External API =====
// 	start := time.Now()
// 	customer, respHeader, err := uc.connectorRepo.GetCustInfoByAeonID(ctx, model.ConnectorCustomerInfoRequest{AeonID: req.AeonID, CustID: req.CustID})
// 	elapsed := time.Since(start).Milliseconds()

// 	// ===== Prepare API Log =====
// 	apiLog := &model.ApiLogs{
// 		RequestID:    requestID,
// 		ServiceName:  "Get Customer Info",
// 		Endpoint:     "https://connectorapi.aeonth.com/Api/Common/GetCustomerInfo",
// 		ReqDatetime:  start,
// 		TimeUsage:    int(elapsed),
// 		RespDatetime: time.Now(),
// 		ReqHeader:    reqHeader,
// 	}

// 	// ===== Marshal Request Body =====
// 	if reqJson, err := json.Marshal(req); err == nil {
// 		apiLog.ReqMessage = string(reqJson)
// 	} else {
// 		log.Println("[UseCase] Failed to marshal request body:", err)
// 		apiLog.ReqMessage = "<failed to marshal>"
// 	}

// 	// ===== Handle Response =====
// 	if err != nil {
// 		log.Println("[UseCase] => Error getting customer:", err)

// 		// Set response error message and status
// 		apiLog.StatusCode = http.StatusGatewayTimeout
// 		apiLog.RespMessage = err.Error()
// 		apiLog.RespHeader = "{}"

// 		// Return error with detail
// 		details := map[string]string{
// 			"connector_api": "Connection issue from System-i",
// 		}
// 		lib.HandleErrorContext(ctx, lib.GatewayTimeout.WithDetails(details))

// 	} else {
// 		// Marshal response body
// 		if respJson, err := json.Marshal(customer); err == nil {
// 			apiLog.StatusCode = http.StatusOK
// 			apiLog.RespMessage = string(respJson)
// 		} else {
// 			log.Println("[UseCase] Failed to marshal response:", err)
// 			apiLog.StatusCode = http.StatusInternalServerError
// 			apiLog.RespMessage = "Failed to marshal response"
// 		}

// 		// Flatten and save response header
// 		if respHeader != nil {
// 			if headerJson, err := json.Marshal(utils.FlattenHeader(respHeader)); err == nil {
// 				apiLog.RespHeader = string(headerJson)
// 			} else {
// 				apiLog.RespHeader = "{}"
// 			}
// 		} else {
// 			apiLog.RespHeader = "{}"
// 		}
// 	}

// 	// ===== Save API Log =====
// 	log.Println("[UseCase] About to save API log...")
// 	if saveErr := uc.repo.SaveApiLog(ctx, apiLog); saveErr != nil {
// 		log.Println("[UseCase] Failed to save API log:", saveErr)
// 	} else {
// 		log.Println("[UseCase] API log saved successfully")
// 	}

// 	return customer, err
// }

func (uc *DashboardUseCase) CustInfo(ctx context.Context, ginCtx *gin.Context, aeonID string) (*model.GetCustInfoResponse, error) {
	log.Printf("[UseCase] => Called with request aeon_id: %+v\n", aeonID)

	// ===== Step 1: Lookup cust_id จาก d_mobile_app_daily =====
	custID, err := uc.repo.GetCustIDByAeonID(ctx, aeonID)
	if err != nil {
		log.Println("[UseCase] => Failed to map aeon_id to cust_id:", err)
		lib.HandleErrorContext(ctx, lib.NotFound.WithDetails("customer not found"))
		return nil, err
	}

	log.Printf("[UseCase] => Mapped aeon_id %s to cust_id %s\n", aeonID, custID)

	// ===== Prepare request header for log =====
	var reqHeader string
	if ginCtx != nil && ginCtx.Request != nil {
		if headerBytes, err := json.Marshal(utils.FlattenHeader(ginCtx.Request.Header)); err == nil {
			reqHeader = string(headerBytes)
		}
	}

	// ===== Get Request ID =====
	requestID, _ := ctx.Value(utils.CtxKeyRequestID).(string)
	if requestID == "" && ginCtx != nil {
		requestID = ginCtx.GetHeader("X-Request-ID")
	}
	if requestID == "" {
		requestID = uuid.New().String()
	}

	// ===== Step 2: Create payload using cust_id =====
	payload := model.ConnectorCustomerInfoRequest{
		UserRef: custID,
		Mode:    "F",
	}

	// ===== Step 3: Call third-party API =====
	start := time.Now()
	customer, respHeader, err := uc.connectorRepo.GetCustInfoByAeonID(ctx, payload)
	elapsed := time.Since(start).Milliseconds()

	// ===== Step 4: Logging and error handling as before =====
	apiLog := &model.ApiLogs{
		RequestID:    requestID,
		ServiceName:  "Get Customer Info",
		Endpoint:     "https://connectorapi.aeonth.com/Api/Common/GetCustomerInfo",
		ReqDatetime:  start,
		TimeUsage:    int(elapsed),
		RespDatetime: time.Now(),
		ReqHeader:    reqHeader,
	}

	if reqJson, err := json.Marshal(payload); err == nil {
		apiLog.ReqMessage = string(reqJson)
	} else {
		apiLog.ReqMessage = "<failed to marshal>"
	}

	if err != nil {
		apiLog.StatusCode = http.StatusGatewayTimeout
		apiLog.RespMessage = err.Error()
		apiLog.RespHeader = "{}"

		log.Println("[UseCase] => Error getting customer:", err)
		lib.HandleErrorContext(ctx, lib.GatewayTimeout.WithDetails(map[string]string{
			"connector_api": "Connection issue from System-i",
		}))
	} else {
		if respJson, err := json.Marshal(customer); err == nil {
			apiLog.StatusCode = http.StatusOK
			apiLog.RespMessage = string(respJson)
		} else {
			apiLog.StatusCode = http.StatusInternalServerError
			apiLog.RespMessage = "Failed to marshal response"
		}

		if respHeader != nil {
			if headerJson, err := json.Marshal(utils.FlattenHeader(respHeader)); err == nil {
				apiLog.RespHeader = string(headerJson)
			}
		}
	}

	log.Println("[UseCase] About to save API log...")
	if saveErr := uc.repo.SaveApiLog(ctx, apiLog); saveErr != nil {
		log.Println("[UseCase] Failed to save API log:", saveErr)
	} else {
		log.Println("[UseCase] API log saved successfully")
	}

	return customer, err
}

func (uc *DashboardUseCase) CustProfile(ctx context.Context, ginCtx *gin.Context, aeonID string) (*model.GetCustProfileResponse, error) {
	log.Printf("[UseCase] => Called CustProfile with aeonID: %s\n", aeonID)

	// ===== Extract Request Header for logging =====
	var reqHeader string
	if ginCtx != nil && ginCtx.Request != nil {
		if headerBytes, err := json.Marshal(utils.FlattenHeader(ginCtx.Request.Header)); err == nil {
			reqHeader = string(headerBytes)
		}
	}

	// ===== Get Request ID from Context or Header =====
	requestID, _ := ctx.Value(utils.CtxKeyRequestID).(string)
	if requestID == "" && ginCtx != nil {
		requestID = ginCtx.GetHeader("X-Request-ID")
	}
	if requestID == "" {
		requestID = uuid.New().String()
	}

	start := time.Now()
	custProfile, respHeader, err := uc.tdRepo.GetCustProfileByAeonID(ctx, aeonID)
	elapsed := time.Since(start).Milliseconds()

	// Prepare log struct (adjust if model.APILog supports these fields)
	apiLog := &model.ApiLogs{
		RequestID:    requestID,
		ServiceName:  "Get Customer Profile",
		Endpoint:     "https://cdp.in.treasuredata.com/cdp/lookup/collect/profiles?version=2&token=",
		ReqDatetime:  start,
		TimeUsage:    int(elapsed),
		RespDatetime: time.Now(),
		ReqHeader:    reqHeader,
		ReqMessage:   aeonID,
	}

	if err != nil {
		log.Printf("[UseCase] => Error CustProfile: %v\n", err)
		apiLog.StatusCode = http.StatusGatewayTimeout
		apiLog.RespMessage = err.Error()
		apiLog.RespHeader = "{}"

		details := map[string]string{
			"cdp_api": "Connection issue from TD",
		}
		lib.HandleErrorContext(ctx, lib.GatewayTimeout.WithDetails(details))
	} else {
		if respJson, err := json.Marshal(custProfile); err == nil {
			apiLog.StatusCode = http.StatusOK
			apiLog.RespMessage = string(respJson)
		} else {
			log.Println("[UseCase] Failed to marshal profile response:", err)
			apiLog.StatusCode = http.StatusInternalServerError
			apiLog.RespMessage = "Failed to marshal response"
		}

		apiLog.RespHeader = "{}" // ปรับถ้าต้องการเก็บ header จริงๆ
	}

	// Flatten and save response header
	if respHeader != nil {
		if headerJson, err := json.Marshal(utils.FlattenHeader(respHeader)); err == nil {
			apiLog.RespHeader = string(headerJson)
		} else {
			apiLog.RespHeader = "{}"
		}
	} else {
		apiLog.RespHeader = "{}"
	}

	log.Println("[UseCase] About to save API log...")
	if saveErr := uc.repo.SaveApiLog(ctx, apiLog); saveErr != nil {
		log.Println("[UseCase] Failed to save API log:", saveErr)
	} else {
		log.Println("[UseCase] API log saved successfully")
	}

	return custProfile, err
}

func (uc *DashboardUseCase) CustSegment(ctx context.Context, ginCtx *gin.Context, aeonID string) (*model.GetCustSegmentResponse, error) {
	log.Printf("[UseCase] => Called CustSegment with aeonID: %s\n", aeonID)

	var reqHeader string
	if ginCtx != nil && ginCtx.Request != nil {
		if headerBytes, err := json.Marshal(utils.FlattenHeader(ginCtx.Request.Header)); err == nil {
			reqHeader = string(headerBytes)
		}
	}

	requestID, _ := ctx.Value(utils.CtxKeyRequestID).(string)
	if requestID == "" && ginCtx != nil {
		requestID = ginCtx.GetHeader("X-Request-ID")
	}
	if requestID == "" {
		requestID = uuid.New().String()
	}

	start := time.Now()
	custSegment, respHeader, err := uc.tdRepo.GetCustSegmentByAeonID(ctx, aeonID)
	elapsed := time.Since(start).Milliseconds()

	apiLog := &model.ApiLogs{
		RequestID:    requestID,
		ServiceName:  "Get Customer Segment",
		Endpoint:     "https://cdp.in.treasuredata.com/cdp/lookup/collect/segments?version=2&token=",
		ReqDatetime:  start,
		TimeUsage:    int(elapsed),
		RespDatetime: time.Now(),
		ReqHeader:    reqHeader,
		ReqMessage:   aeonID,
	}

	if err != nil {
		log.Printf("[UseCase] => Error CustSegment: %v\n", err)
		apiLog.StatusCode = http.StatusGatewayTimeout
		apiLog.RespMessage = err.Error()
		apiLog.RespHeader = "{}"

		details := map[string]string{
			"cdp_api": "Connection issue from TD",
		}
		lib.HandleErrorContext(ctx, lib.GatewayTimeout.WithDetails(details))
	} else {
		if respJson, err := json.Marshal(custSegment); err == nil {
			apiLog.StatusCode = http.StatusOK
			apiLog.RespMessage = string(respJson)
		} else {
			log.Println("[UseCase] Failed to marshal segment response:", err)
			apiLog.StatusCode = http.StatusInternalServerError
			apiLog.RespMessage = "Failed to marshal response"
		}

		apiLog.RespHeader = "{}"
	}

	// Flatten and save response header
	if respHeader != nil {
		if headerJson, err := json.Marshal(utils.FlattenHeader(respHeader)); err == nil {
			apiLog.RespHeader = string(headerJson)
		} else {
			apiLog.RespHeader = "{}"
		}
	} else {
		apiLog.RespHeader = "{}"
	}

	log.Println("[UseCase] About to save API log...")
	if saveErr := uc.repo.SaveApiLog(ctx, apiLog); saveErr != nil {
		log.Println("[UseCase] Failed to save API log:", saveErr)
	} else {
		log.Println("[UseCase] API log saved successfully")
	}

	return custSegment, err
}

func (uc *DashboardUseCase) CustSuggestion(ctx context.Context, ginCtx *gin.Context, aeonID string) (*model.GetCustSuggestionResponse, error) {
	log.Printf("[UseCase] => Called CustSuggestion with aeonID: %s\n", aeonID)

	var reqHeader string
	if ginCtx != nil && ginCtx.Request != nil {
		if headerBytes, err := json.Marshal(utils.FlattenHeader(ginCtx.Request.Header)); err == nil {
			reqHeader = string(headerBytes)
		}
	}

	requestID, _ := ctx.Value(utils.CtxKeyRequestID).(string)
	if requestID == "" && ginCtx != nil {
		requestID = ginCtx.GetHeader("X-Request-ID")
	}
	if requestID == "" {
		requestID = uuid.New().String()
	}

	start := time.Now()
	custSuggestion, respHeader, err := uc.tdRepo.GetCustSuggestionByAeonID(ctx, aeonID)
	elapsed := time.Since(start).Milliseconds()

	apiLog := &model.ApiLogs{
		RequestID:    requestID,
		ServiceName:  "Get Customer Suggestion",
		Endpoint:     "https://cdp.in.treasuredata.com/cdp/lookup/collect/segments?version=2&token=",
		ReqDatetime:  start,
		TimeUsage:    int(elapsed),
		RespDatetime: time.Now(),
		ReqHeader:    reqHeader,
		ReqMessage:   aeonID,
	}

	if err != nil {
		log.Printf("[UseCase] => Error CustSuggestion: %v\n", err)
		apiLog.StatusCode = http.StatusGatewayTimeout
		apiLog.RespMessage = err.Error()
		apiLog.RespHeader = "{}"

		details := map[string]string{
			"cdp_api": "Connection issue from TD",
		}
		lib.HandleErrorContext(ctx, lib.GatewayTimeout.WithDetails(details))
	} else {
		if respJson, err := json.Marshal(custSuggestion); err == nil {
			apiLog.StatusCode = http.StatusOK
			apiLog.RespMessage = string(respJson)
		} else {
			log.Println("[UseCase] Failed to marshal suggestion response:", err)
			apiLog.StatusCode = http.StatusInternalServerError
			apiLog.RespMessage = "Failed to marshal response"
		}

		apiLog.RespHeader = "{}"
	}

	// Flatten and save response header
	if respHeader != nil {
		if headerJson, err := json.Marshal(utils.FlattenHeader(respHeader)); err == nil {
			apiLog.RespHeader = string(headerJson)
		} else {
			apiLog.RespHeader = "{}"
		}
	} else {
		apiLog.RespHeader = "{}"
	}

	log.Println("[UseCase] About to save API log...")
	if saveErr := uc.repo.SaveApiLog(ctx, apiLog); saveErr != nil {
		log.Println("[UseCase] Failed to save API log:", saveErr)
	} else {
		log.Println("[UseCase] API log saved successfully")
	}

	return custSuggestion, err
}
