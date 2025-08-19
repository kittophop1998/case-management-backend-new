package http

import (
	"case-management/infrastructure/lib"
	"case-management/internal/app/usecase"
	"case-management/internal/domain/model"
	"net/http"

	"github.com/gin-gonic/gin"
)

type CaseHandler struct {
	UseCase usecase.CaseUseCase
}

func (h *CaseHandler) CreateCase(ctx *gin.Context) {
	caseData := &model.CreateCaseRequest{}
	if err := ctx.ShouldBindJSON(caseData); err != nil {
		lib.HandleError(ctx, lib.BadRequest.WithDetails(err.Error()))
		return
	}

	caseID, err := h.UseCase.CreateCase(ctx, caseData)
	if err != nil {
		ctx.JSON(500, gin.H{"error": "Failed to create case"})
		return
	}

	lib.HandleResponse(ctx, http.StatusCreated, gin.H{"caseId": caseID})
}

func (h *CaseHandler) GetAllCases(ctx *gin.Context) {
	cases, err := h.UseCase.GetAllCases(ctx)
	if err != nil {
		ctx.JSON(500, gin.H{"error": "Failed to retrieve cases"})
		return
	}
	ctx.JSON(200, cases)
}

func (h *CaseHandler) GetCaseByID(ctx *gin.Context) {
	caseID := ctx.Param("id")
	caseData, err := h.UseCase.GetCaseByID(ctx, caseID)
	if err != nil {
		ctx.JSON(500, gin.H{"error": "Failed to retrieve case"})
		return
	}
	ctx.JSON(200, caseData)
}

func (h *CaseHandler) GetAllDisposition(ctx *gin.Context) {
	limit, err := getLimit(ctx)
	if err != nil {
		ctx.JSON(400, gin.H{"error": "Invalid limit parameter"})
		return
	}

	page, err := getPage(ctx)
	if err != nil {
		ctx.JSON(400, gin.H{"error": "Invalid page parameter"})
		return
	}

	keyword := ctx.Query("keyword")

	mains, total, err := h.UseCase.GetAllDisposition(ctx, page, limit, keyword)
	if err != nil {
		ctx.JSON(500, gin.H{"error": err.Error()})
		return
	}

	if mains == nil {
		mains = []model.DispositionMain{}
	}

	lib.HandlePaginatedResponse(ctx, page, limit, total, mains)
}
