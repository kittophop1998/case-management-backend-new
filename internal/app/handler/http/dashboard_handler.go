package http

import (
	"case-management/internal/app/usecase"
	"net/http"

	"github.com/gin-gonic/gin"
)

type DashboardHandler struct {
	UseCase usecase.DashboardUseCase
}

func (h *DashboardHandler) GetCustomerInfo(ctx *gin.Context) {
	id := ctx.Param("id")
	customer, err := h.UseCase.CustomerInfo(ctx, id)
	if err != nil {
		HandleError(ctx, err)
		return
	}
	ctx.JSON(http.StatusOK, customer)
}
