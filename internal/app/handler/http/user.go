package http

import (
	"case-management/internal/app/usecase"
	"case-management/internal/domain/model"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type UserHandler struct {
	useCase usecase.UserUseCase
}

func (h *UserHandler) GetAllUsers(c *gin.Context) {
	users, err := h.useCase.GetAll(c)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, users)
}

func (h *UserHandler) GetUserByID(c *gin.Context) {
	id := c.Param("id")
	uid, err := uuid.Parse(id)
	if err != nil {
		c.JSON(400, gin.H{"error": "Invalid user ID"})
		return
	}

	user, err := h.useCase.GetById(c, uid)
	if err != nil {
		c.JSON(404, gin.H{"error": "User not found"})
		return
	}

	c.JSON(200, user)
}

func (h *UserHandler) UpdateUserByID(c *gin.Context) {
	var input model.CreateUpdateUserRequest

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	idParam := c.Param("id")
	userID, err := uuid.Parse(idParam)
	if err != nil {
		c.JSON(400, gin.H{"error": "Invalid user ID"})
		return
	}

	err = h.useCase.Update(c, userID, input)
	if err != nil {
		c.JSON(400, gin.H{"error": "Invalid user ID"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "user updated successfully"})
}
