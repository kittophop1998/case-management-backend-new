package lib

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

type ResponseError struct {
	Code    int          `json:"code"`
	Message MessageError `json:"message"`
	Detail  interface{}  `json:"detail"`
}

type MessageError struct {
	Th string `json:"th"`
	En string `json:"en"`
}

func (e *ResponseError) Error() string {
	return e.Message.Th
}

func NewAppError(code int, message MessageError, detail interface{}) *ResponseError {
	return &ResponseError{
		Code:    code,
		Message: message,
		Detail:  detail,
	}
}

var (
	BadRequest = NewAppError(http.StatusBadRequest, MessageError{
		Th: "คำขอไม่ถูกต้อง",
		En: "Bad Request",
	}, nil)
	InternalServer = NewAppError(http.StatusInternalServerError, MessageError{
		Th: "เกิดข้อผิดพลาดภายในเซิร์ฟเวอร์",
		En: "Internal Server Error",
	}, nil)
	GatewayTimeout = NewAppError(http.StatusGatewayTimeout, MessageError{
		Th: "หมดเวลาการเชื่อมต่อ",
		En: "Gateway Timeout",
	}, nil)
)

// HandleError is a centralized error handler for Gin
func HandleError(ctx *gin.Context, err error) {
	var apiError *ResponseError
	if errors.As(err, &apiError) {
		ctx.AbortWithStatusJSON(apiError.Code, apiError)
		return
	}
}

func (e *ResponseError) WithDetails(details interface{}) *ResponseError {
	return &ResponseError{
		Code:    e.Code,
		Message: e.Message,
		Detail:  details,
	}
}

func (e *ResponseError) WithMessage(message MessageError) *ResponseError {
	return &ResponseError{
		Code:    e.Code,
		Message: message,
		Detail:  e.Detail,
	}
}
