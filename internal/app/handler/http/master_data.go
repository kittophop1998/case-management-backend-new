package http

import (
	"case-management/internal/app/usecase"
	"net/http"

	"github.com/gin-gonic/gin"
)

type MasterDataHandler struct {
	UseCase usecase.MasterDataUseCase
}

func (h *MasterDataHandler) GetAllLookups(ctx *gin.Context) {
	data, err := h.UseCase.GetAllLookups(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": data})
}
