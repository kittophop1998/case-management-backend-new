package http

import (
	"case-management/internal/app/usecase"
	"case-management/internal/domain/model"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type UserHandler struct {
	UseCase usecase.UserUseCase
}

func (h *UserHandler) GetAllUsers(ctx *gin.Context) {
	users, err := h.UseCase.GetAll(ctx)
	if err != nil {
		ctx.JSON(500, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(200, users)
}

func (h *UserHandler) GetUserByID(ctx *gin.Context) {
	id := ctx.Param("id")
	userId, err := uuid.Parse(id)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	user, err := h.UseCase.GetById(ctx, userId)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	ctx.JSON(http.StatusOK, user)
}

func (h *UserHandler) UpdateUserByID(ctx *gin.Context) {
	var input model.CreateUpdateUserRequest
	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(500, gin.H{"error": err.Error()})
		return
	}

	id := ctx.Param("id")
	userId, err := uuid.Parse(id)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	err = h.UseCase.UpdateUserById(ctx, userId, input)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Can't update user"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "user updated successfully"})
}
