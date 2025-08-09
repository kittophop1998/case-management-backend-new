package database

import (
	"case-management/internal/domain/model"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type LogPg struct {
	db *gorm.DB
}

func NewLogPg(db *gorm.DB) *LogPg {
	return &LogPg{db: db}
}

func (l *LogPg) SaveApiLog(log *model.ApiLogs) error {
	return l.db.Create(log).Error
}

func (l *LogPg) GetAllApiLogs(ctx *gin.Context) ([]*model.ApiLogs, error) {
	var logs []*model.ApiLogs

	if err := l.db.WithContext(ctx).Order("created_at desc").Find(&logs).Error; err != nil {
		return nil, err
	}

	return logs, nil
}
