package http

import (
	"errors"
	"time"

	"case-management/pkg/api_errors"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func HandleError(ctx *gin.Context, err error) {
	var appErr *api_errors.AppError
	//var requestID string

	//if id, exists := c.Get(RequestIDKey); exists {
	//	requestID = id.(string)
	//}

	if errors.As(err, &appErr) {
		ctx.JSON(appErr.HTTPStatus, api_errors.APIErrorResponse{
			Errors: []api_errors.ErrorDetail{
				{
					Code:      appErr.Code,
					Message:   appErr.Message,
					Details:   appErr.Details,
					Status:    appErr.HTTPStatus,
					Timestamp: time.Now().Format(time.RFC3339Nano),
				},
			},
		})
		return
	}

	if errors.Is(err, gorm.ErrRecordNotFound) {
		HandleError(ctx, api_errors.ErrNotFound)
		return
	}

	unhandledErr := api_errors.ErrInternalServer
	ctx.JSON(unhandledErr.HTTPStatus, api_errors.APIErrorResponse{
		Errors: []api_errors.ErrorDetail{
			{
				Code:      unhandledErr.Code,
				Message:   unhandledErr.Message,
				Details:   map[string]string{"original_error": err.Error()},
				Status:    unhandledErr.HTTPStatus,
				Timestamp: time.Now().Format(time.RFC3339Nano),
			},
		},
	})
}
