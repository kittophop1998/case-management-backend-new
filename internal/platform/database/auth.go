package database

import (
	"case-management/internal/domain/model"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type AuthPg struct {
	db *gorm.DB
}

func NewAuthPg(db *gorm.DB) *AuthPg {
	return &AuthPg{db: db}
}

func (repo *AuthPg) SaveAccessLog(ctx *gin.Context, accessLog *model.AccessLogs) error {
	if err := repo.db.WithContext(ctx).Create(accessLog).Error; err != nil {
		return err
	}
	return nil
}
