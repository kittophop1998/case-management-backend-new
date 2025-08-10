package http

import (
	"case-management/internal/app/usecase"
	"case-management/internal/domain/model"

	"github.com/gin-gonic/gin"
)

type CaseHandler struct {
	UseCase usecase.CaseUseCase
}

func (h *CaseHandler) CreateCase(ctx *gin.Context) {
	caseData := &model.CreateCaseRequest{}
	if err := ctx.ShouldBindJSON(caseData); err != nil {
		ctx.JSON(400, gin.H{"error": "Invalid case data"})
		return
	}

	caseID, err := h.UseCase.CreateCase(ctx, caseData)
	if err != nil {
		ctx.JSON(500, gin.H{"error": "Failed to create case"})
		return
	}
	ctx.JSON(201, gin.H{"case_id": caseID})
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
