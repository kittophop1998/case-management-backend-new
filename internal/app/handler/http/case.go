package http

import (
	"case-management/infrastructure/lib"
	"case-management/internal/app/usecase"
	"case-management/internal/domain/model"
	"case-management/utils"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type CaseHandler struct {
	UseCase usecase.CaseUseCase
}

func (h *CaseHandler) CreateCaseInquiry(ctx *gin.Context) {
	userId := ctx.GetString("userId")
	createdByID, err := uuid.Parse(userId)
	if err != nil {
		lib.HandleError(ctx, lib.InternalServer.WithDetails(err.Error()))
		return
	}

	caseReq := &model.CreateCaseInquiryRequest{}
	if err := ctx.ShouldBindJSON(caseReq); err != nil {
		lib.HandleError(ctx, lib.BadRequest.WithDetails(err.Error()))
		return
	}

	caseID, err := h.UseCase.CreateCaseInquiry(ctx, createdByID, caseReq)
	if err != nil {
		lib.HandleError(ctx, lib.InternalServer.WithDetails(err.Error()))
		return
	}

	lib.HandleResponse(ctx, http.StatusCreated, gin.H{"caseId": caseID})
}

func (h *CaseHandler) GetAllCases(ctx *gin.Context) {
	p := utils.GetPagination(ctx)

	// TODO: Implement filter category
	// category := ctx.Query("category")

	cases, total, err := h.UseCase.GetAllCases(ctx, p.Page, p.Limit)
	if err != nil {
		lib.HandleError(ctx, lib.InternalServer.WithDetails(err.Error()))
		return
	}

	lib.HandlePaginatedResponse(ctx, p.Page, p.Limit, total, cases)
}

func (h *CaseHandler) GetCaseByID(ctx *gin.Context) {
	caseID := ctx.Param("id")
	caseData, err := h.UseCase.GetCaseByID(ctx, caseID)
	if err != nil {
		ctx.JSON(500, gin.H{"error": "Failed to retrieve case"})
		return
	}
	lib.HandleResponse(ctx, http.StatusOK, caseData)
}

func (h *CaseHandler) GetAllDisposition(ctx *gin.Context) {
	mains, err := h.UseCase.GetAllDisposition(ctx)
	if err != nil {
		lib.HandleError(ctx, lib.InternalServer.WithDetails(err.Error()))
		return
	}

	if mains == nil {
		mains = []model.DispositionItem{}
	}

	lib.HandleResponse(ctx, http.StatusOK, mains)
}
