package lib

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type ErrorDetail struct {
	TimeStamp  string       `json:"timestamp"`
	Path       string       `json:"path"`
	StatusCode string       `json:"statusCode"`
	Code       int          `json:"code"`
	Message    MessageError `json:"message"`
	Detail     interface{}  `json:"detail"`
}

type MessageError struct {
	Th string `json:"th"`
	En string `json:"en"`
}

func (e *ErrorDetail) Error() string {
	return e.Message.Th
}

func NewAppError(code int, message MessageError, detail interface{}) *ErrorDetail {
	return &ErrorDetail{
		StatusCode: http.StatusText(code),
		Code:       code,
		Message:    message,
		Detail:     detail,
	}
}

// Common predefined errors
var (
	BadRequest     = NewAppError(http.StatusBadRequest, MessageError{"คำขอไม่ถูกต้อง", "Bad Request"}, nil)
	InternalServer = NewAppError(http.StatusInternalServerError, MessageError{"เกิดข้อผิดพลาดภายในเซิร์ฟเวอร์", "Internal Server Error"}, nil)
	GatewayTimeout = NewAppError(http.StatusGatewayTimeout, MessageError{"หมดเวลาการเชื่อมต่อ", "Gateway Timeout"}, nil)
	Unauthorized   = NewAppError(http.StatusUnauthorized, MessageError{"ไม่อนุญาต", "Unauthorized"}, nil)
	NotFound       = NewAppError(http.StatusNotFound, MessageError{"ไม่พบข้อมูล", "Not Found"}, nil)
	CannotUpdate   = NewAppError(http.StatusConflict, MessageError{"ไม่สามารถอัปเดตข้อมูลได้", "Cannot Update"}, nil)
	CannotCreate   = NewAppError(http.StatusConflict, MessageError{"ไม่สามารถสร้างข้อมูลได้", "Cannot Create"}, nil)
)

type ApiErrorResponse struct {
	Error ErrorDetail `json:"error"`
}

func timestampNow() string {
	return time.Now().UTC().Format(time.RFC3339Nano)
}

// HandleError is a centralized error handler for Gin
func HandleError(ctx *gin.Context, err error) {
	var apiError *ErrorDetail
	if errors.As(err, &apiError) {
		ctx.AbortWithStatusJSON(apiError.Code, ApiErrorResponse{
			Error: ErrorDetail{
				TimeStamp:  timestampNow(),
				Path:       ctx.FullPath(),
				StatusCode: apiError.StatusCode,
				Code:       apiError.Code,
				Message:    apiError.Message,
				Detail:     apiError.Detail,
			},
		})
		return
	}

	// fallback for unexpected errors
	ctx.AbortWithStatusJSON(http.StatusInternalServerError, ApiErrorResponse{
		Error: ErrorDetail{
			TimeStamp:  timestampNow(),
			Path:       ctx.FullPath(),
			StatusCode: http.StatusText(apiError.Code),
			Code:       http.StatusInternalServerError,
			Message:    MessageError{"เกิดข้อผิดพลาดที่ไม่คาดคิด", "Unexpected Error"},
			Detail:     err.Error(),
		},
	})
}

// lib/error_handler.go
func HandleErrorContext(ctx context.Context, err error) *ApiErrorResponse {
	var apiError *ErrorDetail
	if errors.As(err, &apiError) {
		return &ApiErrorResponse{
			Error: ErrorDetail{
				TimeStamp: timestampNow(),
				// Path:       ctx.Value("path").(string),
				StatusCode: apiError.StatusCode,
				Code:       apiError.Code,
				Message:    apiError.Message,
				Detail:     apiError.Detail,
			},
		}
	}

	// fallback
	return &ApiErrorResponse{
		Error: ErrorDetail{
			TimeStamp: timestampNow(),
			// Path:       ctx.Value("path").(string),
			StatusCode: http.StatusText(http.StatusInternalServerError),
			Code:       http.StatusInternalServerError,
			Message:    MessageError{"เกิดข้อผิดพลาดที่ไม่คาดคิด", "Unexpected Error"},
			Detail:     err.Error(),
		},
	}
}

// Chainable methods
func (e *ErrorDetail) WithDetails(details interface{}) *ErrorDetail {
	newErr := *e
	newErr.Detail = details
	return &newErr
}

func (e *ErrorDetail) WithMessage(message MessageError) *ErrorDetail {
	newErr := *e
	newErr.Message = message
	return &newErr
}
