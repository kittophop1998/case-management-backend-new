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
	UseCase          usecase.CaseUseCase
	UpdateCaseByType usecase.UpdateCaseUseCase
}

func (h *CaseHandler) CreateCase(c *gin.Context) {
	userIdRaw, exists := c.Get("userId")
	if !exists {
		lib.HandleError(c, lib.InternalServer.WithDetails("userId not found"))
		return
	}
	createdByID, err := uuid.Parse(userIdRaw.(string))
	if err != nil {
		lib.HandleError(c, lib.InternalServer.WithDetails(err.Error()))
		return
	}

	caseReq := &model.CreateCaseRequest{}
	if err := c.ShouldBindJSON(caseReq); err != nil {
		lib.HandleError(c, lib.BadRequest.WithDetails(err.Error()))
		return
	}

	ctx := c.Request.Context()
	caseID, err := h.UseCase.CreateCase(ctx, createdByID, caseReq)
	if err != nil {
		lib.HandleError(c, lib.InternalServer.WithDetails(err.Error()))
		return
	}

	lib.HandleResponse(c, http.StatusCreated, gin.H{"caseId": caseID})
}

func (h *CaseHandler) GetAllCases(c *gin.Context) {
	ctx := c.Request.Context()

	userIdRaw, exists := c.Get("userId")
	if !exists {
		lib.HandleError(c, lib.InternalServer.WithDetails("userId not found"))
		return
	}
	currentUserID, err := uuid.Parse(userIdRaw.(string))
	if err != nil {
		lib.HandleError(c, lib.InternalServer.WithDetails("invalid userId"))
		return
	}

	p := utils.GetPagination(c)

	var q model.CaseQuery
	if err := c.ShouldBindQuery(&q); err != nil {
		lib.HandleError(c, lib.BadRequest.WithDetails("invalid query params"))
		return
	}

	sort := c.DefaultQuery("sort", "created_at")

	// compose filter
	filter := model.CaseFilter{
		Keyword:  q.Keyword,
		StatusID: q.StatusID,
		QueueID:  q.QueueID,
		Priority: q.Priority,
		Sort:     sort,
	}

	// call usecase
	cases, total, err := h.UseCase.GetAllCases(ctx, p.Page, p.Limit, filter, q.Category, currentUserID)
	if err != nil {
		lib.HandleError(c, lib.InternalServer.WithDetails(err.Error()))
		return
	}

	lib.HandlePaginatedResponse(c, p.Page, p.Limit, total, cases)
}

func (h *CaseHandler) GetCaseByID(c *gin.Context) {
	ctx := c.Request.Context()

	id := c.Param("id")
	caseID, err := uuid.Parse(id)
	if err != nil {
		lib.HandleError(c, lib.BadRequest.WithDetails(err.Error()))
		return
	}

	caseData, err := h.UseCase.GetCaseByID(ctx, caseID)
	if err != nil {
		lib.HandleError(c, lib.InternalServer.WithDetails(err.Error()))
		return
	}

	lib.HandleResponse(c, http.StatusOK, caseData)
}

func (h *CaseHandler) GetAllDisposition(c *gin.Context) {
	ctx := c.Request.Context()

	mains, err := h.UseCase.GetAllDisposition(ctx)
	if err != nil {
		lib.HandleError(c, lib.InternalServer.WithDetails(err.Error()))
		return
	}

	if mains == nil {
		mains = []model.DispositionItem{}
	}

	lib.HandleResponse(c, http.StatusOK, mains)
}

func (h *CaseHandler) UpdateCaseByID(c *gin.Context) {
	ctx := c.Request.Context()

	caseID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		lib.HandleError(c, lib.BadRequest.WithDetails("Invalid case ID"))
		return
	}

	var req model.UpdateCaseRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		lib.HandleError(c, lib.BadRequest.WithDetails(err.Error()))
		return
	}

	if req.CaseTypeID == nil {
		lib.HandleError(c, lib.BadRequest.WithDetails("caseTypeId is required"))
		return
	}

	caseTypeID, err := uuid.Parse(*req.CaseTypeID)
	if err != nil {
		lib.HandleError(c, lib.BadRequest.WithDetails("Invalid case type ID"))
		return
	}

	caseType, err := h.UseCase.GetCaseTypeByID(ctx, caseTypeID)
	if err != nil {
		lib.HandleError(c, lib.InternalServer.WithDetails(err.Error()))
		return
	}

	if err := h.UseCase.UpdateCaseDetail(ctx, caseID, &req); err != nil {
		lib.HandleError(c, lib.InternalServer.WithDetails(err.Error()))
		return
	}

	if err := h.UpdateCaseByType.Execute(ctx, caseType.Group, caseID, req.Data); err != nil {
		lib.HandleError(c, lib.InternalServer.WithDetails(err.Error()))
		return
	}

	lib.HandleResponse(c, http.StatusOK, gin.H{"message": "Case updated successfully"})
}

func (h *CaseHandler) AddCaseNote(c *gin.Context) {
	ctx := c.Request.Context()

	userIdRaw, exists := c.Get("userId")
	if !exists {
		lib.HandleError(c, lib.InternalServer.WithDetails("userId not found"))
		return
	}
	createdByID, err := uuid.Parse(userIdRaw.(string))
	if err != nil {
		lib.HandleError(c, lib.InternalServer.WithDetails(err.Error()))
		return
	}

	caseId := c.Param("id")
	caseID, err := uuid.Parse(caseId)
	if err != nil {
		lib.HandleError(c, lib.BadRequest.WithDetails(err.Error()))
		return
	}

	noteReq := &model.CaseNoteRequest{}
	if err := c.ShouldBindJSON(noteReq); err != nil {
		lib.HandleError(c, lib.BadRequest.WithDetails(err.Error()))
		return
	}

	noteID, err := h.UseCase.AddCaseNote(ctx, createdByID, caseID, noteReq)
	if err != nil {
		lib.HandleError(c, lib.InternalServer.WithDetails(err.Error()))
		return
	}

	lib.HandleResponse(c, http.StatusCreated, gin.H{"noteId": noteID})
}

func (h *CaseHandler) GetCaseNotes(c *gin.Context) {
	ctx := c.Request.Context()

	caseIdRaw := c.Param("caseId")
	caseID, err := uuid.Parse(caseIdRaw)
	if err != nil {
		lib.HandleError(c, lib.BadRequest.WithDetails("Invalid case ID"))
		return
	}

	notes, err := h.UseCase.GetCaseNotes(ctx, caseID)
	if err != nil {
		lib.HandleError(c, lib.InternalServer.WithDetails(err.Error()))
		return
	}

	lib.HandleResponse(c, http.StatusOK, notes)
}
