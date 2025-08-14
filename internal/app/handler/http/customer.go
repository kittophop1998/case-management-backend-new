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
	note := &model.CustomerNote{}
	if err := ctx.ShouldBindJSON(note); err != nil {
		lib.HandleError(ctx, lib.BadRequest.WithDetails(err.Error()))
		return
	}

	if err := h.UseCase.CreateCustomerNote(note); err != nil {
		lib.HandleError(ctx, lib.InternalServer.WithDetails(err.Error()))
		return
	}
	lib.HandleResponse(ctx, http.StatusCreated, gin.H{"message": "Customer note created successfully"})
}

func (h *CustomerHandler) GetAllCustomerNotes(ctx *gin.Context) {
	customerID := ctx.Param("customerId")
	notes, err := h.UseCase.GetAllCustomerNotes(customerID)
	if err != nil {
		lib.HandleError(ctx, lib.InternalServer.WithDetails(err.Error()))
		return
	}

	lib.HandleResponse(ctx, http.StatusOK, notes)
}
