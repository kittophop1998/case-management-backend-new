package http

import (
	"case-management/internal/app/usecase"
	"case-management/internal/domain/model"
	"net/http"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	UseCase usecase.AuthUseCase
}

func (h *AuthHandler) Login(ctx *gin.Context) {
	var req model.LoginRequest

	// Validate body
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request body",
			"details": err.Error(),
		})
		return
	}

	if req.Username == "" || req.Password == "" {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error":   "Missing credentials",
			"details": "Please provide username and password",
		})
		return
	}

	// Main login use case
	resp, err := h.UseCase.Login(ctx, req)
	success := err == nil

	// Log access (even on failure)
	_ = h.UseCase.SaveAccessLog(ctx, req.Username, success)

	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error":   "Login failed",
			"details": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, resp)
}

func (h *AuthHandler) Logout(ctx *gin.Context) {
	if err := h.UseCase.Logout(ctx); err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error":   "Logout failed",
			"details": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Logged out successfully"})
}

func (h *AuthHandler) Profile(ctx *gin.Context) {
	user, err := h.UseCase.GetProfile(ctx)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to fetch profile",
			"details": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, user)
}
