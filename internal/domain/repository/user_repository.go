package repository

import (
	"case-management/internal/domain/model"

	"github.com/gin-gonic/gin"
)

type UserRepository interface {
	FindAll(c *gin.Context) ([]model.User, error)
	// FindByID(c *gin.Context, id string) (*model.User, error)
}
