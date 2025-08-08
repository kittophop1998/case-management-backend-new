package http

import (
	"case-management/internal/app/usecase"
	"net/http"

	"github.com/gin-gonic/gin"
)

type MasterDataHandler struct {
	UseCase usecase.MasterDataUseCase
}

func (h *MasterDataHandler) GetAllLookups(c *gin.Context) {
	data, err := h.UseCase.GetAllLookups(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": data})
}
