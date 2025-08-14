package http

import (
	"case-management/infrastructure/lib"
	"case-management/internal/app/usecase"
	"case-management/internal/domain/model"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type AuthHandler struct {
	UseCase     usecase.AuthUseCase
	UserUseCase usecase.UserUseCase
}

func (h *AuthHandler) Login(ctx *gin.Context) {
	var req model.LoginRequest

	// Validate body
	if err := ctx.ShouldBindJSON(&req); err != nil {
		lib.HandleError(ctx, lib.BadRequest.WithDetails(err.Error()))
		return
	}

	if req.Username == "" || req.Password == "" {
		lib.HandleError(ctx, lib.BadRequest.WithDetails("Missing credentials"))
		return
	}

	// Main login use case
	resp, err := h.UseCase.Login(ctx, req)
	success := err == nil
	if !success {
		lib.HandleError(ctx, lib.Unauthorized.WithDetails(err.Error()))
		return
	}

	// Log access (even on failure)
	_ = h.UseCase.SaveAccessLog(ctx, req.Username, success)
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
	id, exists := ctx.Get("userId")
	if !exists {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	idStr, ok := id.(string)
	if !ok {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Error id is not string"})
		return
	}

	userId, err := uuid.Parse(idStr)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Error parsing user ID"})
		return
	}

	user, err := h.UserUseCase.GetById(ctx, userId)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to fetch profile",
			"details": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, user)
}
