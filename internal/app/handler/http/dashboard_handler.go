package http

import (
	"case-management/infrastructure/lib"
	"case-management/internal/app/usecase"
	"net/http"

	"github.com/gin-gonic/gin"
)

type DashboardHandler struct {
	UseCase usecase.DashboardUseCase
}

func (h *DashboardHandler) GetCustInfo(ctx *gin.Context) {
	id := ctx.Param("id")
	customer, err := h.UseCase.CustInfo(ctx, id)
	if err != nil {
		lib.HandleError(ctx, lib.BadRequest.WithDetails(err.Error()))
		return
	}
	ctx.JSON(http.StatusOK, customer)
}

func (h *DashboardHandler) GetCustProfile(ctx *gin.Context) {
	id := ctx.Param("id")
	customer, err := h.UseCase.CustProfile(ctx, id)
	if err != nil {
		lib.HandleError(ctx, lib.BadRequest.WithDetails(err.Error()))
		return
	}
	ctx.JSON(http.StatusOK, customer)
}

func (h *DashboardHandler) GetCustSegment(ctx *gin.Context) {
	id := ctx.Param("id")
	customer, err := h.UseCase.CustSegment(ctx, id)
	if err != nil {
		lib.HandleError(ctx, lib.BadRequest.WithDetails(err.Error()))
		return
	}
	ctx.JSON(http.StatusOK, customer)
}

func (h *DashboardHandler) GetCustSuggestion(ctx *gin.Context) {
	id := ctx.Param("id")
	customer, err := h.UseCase.CustSuggestion(ctx, id)
	if err != nil {
		lib.HandleError(ctx, lib.BadRequest.WithDetails(err.Error()))
		return
	}
	ctx.JSON(http.StatusOK, customer)
}
