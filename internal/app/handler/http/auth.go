package http

import (
	"case-management/internal/app/usecase"
	"case-management/internal/domain/model"
	"net/http"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	useCase usecase.AuthUseCase
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req model.LoginRequest

	// Validate body
	if err := c.ShouldBindJSON(&req); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request body",
			"details": err.Error(),
		})
		return
	}

	if req.Username == "" || req.Password == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error":   "Missing credentials",
			"details": "Please provide username and password",
		})
		return
	}

	// Main login use case
	resp, err := h.useCase.Login(c, req)
	success := err == nil

	// Log access (even on failure)
	_ = h.useCase.SaveAccessLog(c, req.Username, success)

	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error":   "Login failed",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, resp)
}

func (h *AuthHandler) Logout(c *gin.Context) {
	if err := h.useCase.Logout(c); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error":   "Logout failed",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Logged out successfully"})
}

func (h *AuthHandler) Profile(c *gin.Context) {
	user, err := h.useCase.GetProfile(c)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to fetch profile",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, user)
}
