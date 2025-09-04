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

func (h *CaseHandler) CreateCase(ctx *gin.Context) {
	userId := ctx.GetString("userId")
	createdByID, err := uuid.Parse(userId)
	if err != nil {
		lib.HandleError(ctx, lib.InternalServer.WithDetails(err.Error()))
		return
	}

	caseReq := &model.CreateCaseRequest{}
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
	userId := ctx.GetString("userId")
	currID, err := uuid.Parse(userId)
	if err != nil {
		lib.HandleError(ctx, lib.InternalServer.WithDetails(err.Error()))
		return
	}

	p := utils.GetPagination(ctx)
	category := ctx.Query("category")

	cases, total, err := h.UseCase.GetAllCases(ctx, p.Page, p.Limit, category, currID)
	if err != nil {
		lib.HandleError(ctx, lib.InternalServer.WithDetails(err.Error()))
		return
	}

	lib.HandlePaginatedResponse(ctx, p.Page, p.Limit, total, cases)
}

func (h *CaseHandler) GetCaseByID(ctx *gin.Context) {
	id := ctx.Param("id")
	caseID, err := uuid.Parse(id)
	if err != nil {
		lib.HandleError(ctx, lib.BadRequest.WithDetails(err.Error()))
		return
	}

	caseData, err := h.UseCase.GetCaseByID(ctx, caseID)
	if err != nil {
		lib.HandleError(ctx, lib.InternalServer.WithDetails(err.Error()))
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

func (h *CaseHandler) AddCaseNote(ctx *gin.Context) {
	userId := ctx.GetString("userId")
	createdByID, err := uuid.Parse(userId)
	if err != nil {
		lib.HandleError(ctx, lib.InternalServer.WithDetails(err.Error()))
		return
	}

	caseId := ctx.Param("id")
	caseID, err := uuid.Parse(caseId)
	if err != nil {
		lib.HandleError(ctx, lib.BadRequest.WithDetails(err.Error()))
		return
	}

	noteReq := &model.CaseNoteRequest{}
	if err := ctx.ShouldBindJSON(noteReq); err != nil {
		lib.HandleError(ctx, lib.BadRequest.WithDetails(err.Error()))
		return
	}

	noteID, err := h.UseCase.AddCaseNote(ctx, createdByID, caseID, noteReq)
	if err != nil {
		lib.HandleError(ctx, lib.InternalServer.WithDetails(err.Error()))
		return
	}

	lib.HandleResponse(ctx, http.StatusCreated, gin.H{"noteId": noteID})
}
