package api_errors

import (
	"net/http"
)

type APIErrorResponse struct {
	Errors []ErrorDetail `json:"error"`
}

type ErrorDetail struct {
	Code      string      `json:"code"`
	Message   string      `json:"message"`
	Details   interface{} `json:"details,omitempty"`
	Status    int         `json:"status"`
	Timestamp string      `json:"timestamp"`
	// RequestID string      `json:"request_id"`
}

type AppError struct {
	Code       string
	Message    string
	HTTPStatus int
	Details    interface{}
}

func (e *AppError) Error() string {
	return e.Message
}

func NewAppError(code, message string, status int, details interface{}) *AppError {
	return &AppError{
		Code:       code,
		Message:    message,
		HTTPStatus: status,
		Details:    details,
	}
}

var (
	ErrBadRequest         = NewAppError("BAD_REQUEST", "ข้อมูลที่ส่งมาไม่ถูกต้องตามที่ระบบรองรับ", http.StatusBadRequest, nil)
	ErrRequiredParam      = NewAppError("REQUIRED_PARAMETER", "ไม่สามารถดึงข้อมูลได้ เนื่องจากระบบขาดข้อมูลบางส่วน", 400, nil)
	ErrFilterRequired     = NewAppError("FILTER_REQUIRED", "โปรดเลือกอย่างน้อย 1 เงื่อนไขเพื่อค้นหา", 400, nil)
	ErrNotFound           = NewAppError("NOT_FOUND", "ไม่พบข้อมูลที่ค้นหา", http.StatusNotFound, nil)
	ErrInternalServer     = NewAppError("INTERNAL_SERVER_ERROR", "เกิดข้อผิดพลาดในระบบ กรุณาลองใหม่ภายหลัง", http.StatusInternalServerError, nil)
	ErrServiceUnavailable = NewAppError("SERVICE_UNAVAILABLE", "The service is temporarily unavailable or in maintenance", http.StatusServiceUnavailable, nil)
	ErrGatewayTimeout     = NewAppError("NO_RESPONSE", "No response from an upstream service", http.StatusGatewayTimeout, nil)
)

func MapError(status int) *AppError {
	if status == 500 {
		return ErrInternalServer
	}

	return ErrInternalServer
}
