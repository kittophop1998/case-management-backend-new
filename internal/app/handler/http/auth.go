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

	if err := ctx.ShouldBindJSON(&req); err != nil {
		lib.HandleError(ctx, lib.BadRequest.WithDetails(err.Error()))
		return
	}

	if req.Username == "" || req.Password == "" {
		lib.HandleError(ctx, lib.BadRequest.WithDetails("Missing credentials"))
		return
	}

	resp, err := h.UseCase.Login(ctx, req)
	success := err == nil
	if !success {
		lib.HandleError(ctx, lib.Unauthorized.WithDetails(err.Error()))
		return
	}

	lib.HandleResponse(ctx, http.StatusOK, resp)
}

func (h *AuthHandler) Logout(ctx *gin.Context) {
	if err := h.UseCase.Logout(ctx); err != nil {
		lib.HandleError(ctx, lib.InternalServer.WithDetails(err.Error()))
		return
	}

	lib.HandleResponse(ctx, http.StatusOK, gin.H{"message": "Logged out successfully"})
}

func (h *AuthHandler) Profile(ctx *gin.Context) {
	userIdRaw, exists := ctx.Get("userId")
	if !exists {
		lib.HandleError(ctx, lib.InternalServer.WithDetails("userId not found"))
		return
	}

	userId, err := uuid.Parse(userIdRaw.(string))
	if err != nil {
		lib.HandleError(ctx, lib.BadRequest.WithDetails(err.Error()))
		return
	}

	user, err := h.UserUseCase.GetProfile(ctx, userId)
	if err != nil {
		lib.HandleError(ctx, lib.InternalServer.WithDetails(err.Error()))
		return
	}

	lib.HandleResponse(ctx, http.StatusOK, user)
}
