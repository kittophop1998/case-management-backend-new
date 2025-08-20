package http

import (
	"case-management/infrastructure/lib"
	"case-management/internal/app/usecase"
	"case-management/internal/domain/model"
	"net/http"

	"github.com/gin-gonic/gin"
)

type CustomerHandler struct {
	UseCase usecase.CustomerUseCase
}

func (h *CustomerHandler) CreateCustomerNote(ctx *gin.Context) {
	username := ctx.GetString("username")

	note := &model.CustomerNote{CreatedBy: username}
	if err := ctx.ShouldBindJSON(note); err != nil {
		lib.HandleError(ctx, lib.BadRequest.WithDetails(err.Error()))
		return
	}

	if err := h.UseCase.CreateCustomerNote(ctx, note); err != nil {
		lib.HandleError(ctx, lib.InternalServer.WithDetails(err.Error()))
		return
	}
	lib.HandleResponse(ctx, http.StatusCreated, gin.H{"message": "Customer note created successfully"})
}

func (h *CustomerHandler) GetAllCustomerNotes(ctx *gin.Context) {
	limit, err := getLimit(ctx)
	if err != nil {
		lib.HandleError(ctx, lib.BadRequest.WithDetails(err.Error()))
		return
	}

	page, err := getPage(ctx)
	if err != nil {
		lib.HandleError(ctx, lib.BadRequest.WithDetails(err.Error()))
		return
	}

	customerID := ctx.Param("customerId")
	notes, total, err := h.UseCase.GetAllCustomerNotes(ctx, customerID, page, limit)
	if err != nil {
		lib.HandleError(ctx, lib.InternalServer.WithDetails(err.Error()))
		return
	}

	lib.HandlePaginatedResponse(ctx, page, limit, total, notes)
}

func (h *CustomerHandler) GetNoteTypes(ctx *gin.Context) {
	noteTypes, err := h.UseCase.GetNoteTypes(ctx)
	if err != nil {
		lib.HandleError(ctx, lib.InternalServer.WithDetails(err.Error()))
		return
	}

	lib.HandleResponse(ctx, http.StatusOK, noteTypes)
}

func (h *CustomerHandler) CountNotes(ctx *gin.Context) {
	customerID := ctx.Param("customerId")
	count, err := h.UseCase.CountNotes(ctx, customerID)
	if err != nil {
		lib.HandleError(ctx, lib.InternalServer.WithDetails(err.Error()))
		return
	}
	lib.HandleResponse(ctx, http.StatusOK, gin.H{"count": count})
}
