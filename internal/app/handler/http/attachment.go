package http

import (
	"case-management/infrastructure/config"
	"case-management/infrastructure/lib"
	"case-management/internal/app/usecase"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type AttachmentHandler struct {
	UseCase usecase.AttachmentUseCase
}

func (h *AttachmentHandler) UploadAttachment(ctx *gin.Context) {
	userIdRaw, exists := ctx.Get("userId")
	if !exists {
		lib.HandleError(ctx, lib.Unauthorized.WithDetails("user_id not found"))
		return
	}
	userId, ok := userIdRaw.(uuid.UUID)
	if !ok {
		lib.HandleError(ctx, lib.InternalServer.WithDetails("invalid user_id format"))
		return
	}

	fileInput, err := ctx.MultipartForm()
	if err != nil {
		lib.HandleError(ctx, lib.BadRequest.WithDetails("Failed to parse form: "+err.Error()))
		return
	}

	files := fileInput.File["file"]
	if len(files) == 0 {
		lib.HandleError(ctx, lib.BadRequest.WithDetails("No file provided"))
		return
	}

	caseIDStr := ctx.Param("case_id")
	caseID, err := uuid.Parse(caseIDStr)
	if err != nil {
		lib.HandleError(ctx, lib.BadRequest.WithDetails("Invalid case_id: "+err.Error()))
		return
	}

	cfg, err := config.Load("")
	if err != nil {
		lib.HandleError(ctx, fmt.Errorf("internal server error"))
	}

	isilonBaseURL := cfg.Isilon.BaseURL
	isilonUser := cfg.Isilon.Username
	isilonPass := cfg.Isilon.Password
	err = h.UseCase.UploadAttachment(ctx, files, caseID, userId, isilonBaseURL, isilonUser, isilonPass)
	if err != nil {
		lib.HandleError(ctx, lib.InternalServer.WithDetails(err.Error()))
		return
	}

	lib.HandleResponse(ctx, http.StatusOK, "File(s) uploaded successfully")
}
