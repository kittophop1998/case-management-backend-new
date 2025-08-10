package http

import (
	"case-management/internal/app/usecase"
	"case-management/internal/domain/model"

	"github.com/gin-gonic/gin"
)

type CustomerHandler struct {
	UseCase usecase.CustomerUseCase
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
	customerID := ctx.Param("customer_id")
	notes, err := h.UseCase.GetAllCustomerNotes(customerID)
	if err != nil {
		ctx.JSON(500, gin.H{"error": "Failed to retrieve customer notes"})
		return
	}
	ctx.JSON(200, gin.H{"data": notes})
}
