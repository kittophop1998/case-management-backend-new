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

func (h *CustomerHandler) SearchByCustomerId(ctx *gin.Context) {
	customerID := ctx.Param("id")
	if customerID == "1234" {
		lib.NewResponse(ctx, http.StatusOK, customerID)
		return
	} else {
		lib.NewResponse(ctx, http.StatusNotFound, "Customer not found")
	}
}

func (h *CustomerHandler) CreateCustomerNote(ctx *gin.Context) {
	note := &model.CustomerNote{}
	if err := ctx.ShouldBindJSON(note); err != nil {
		ctx.JSON(400, gin.H{"error": "Invalid note data"})
		return
	}

	if err := h.UseCase.CreateCustomerNote(note); err != nil {
		ctx.JSON(500, gin.H{"error": "Failed to create customer note"})
		return
	}
	ctx.JSON(201, gin.H{"message": "Customer note created successfully"})
}

func (h *CustomerHandler) GetAllCustomerNotes(ctx *gin.Context) {
	customerID := ctx.Param("customerId")
	notes, err := h.UseCase.GetAllCustomerNotes(customerID)
	if err != nil {
		lib.HandleError(ctx, lib.InternalServer.WithDetails(err.Error()))
		return
	}

	lib.NewResponse(ctx, http.StatusOK, notes)
}
