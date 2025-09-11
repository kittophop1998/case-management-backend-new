package database

import (
	"case-management/internal/domain/model"
	"context"
	"fmt"
	"strings"
	"time"

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

func (l *LogPg) GetAllApiLogs(ctx context.Context, limit, offset int, q *model.APILogQueryParams) ([]*model.ApiLogs, int, error) {
	var logs []*model.ApiLogs

	baseDB := l.db.WithContext(ctx).Model(&model.ApiLogs{})

	if q.RequestID != "" {
		baseDB = baseDB.Where("request_id ILIKE ?", "%"+q.RequestID+"%")
	}
	if q.ServiceName != "" {
		baseDB = baseDB.Where("service_name ILIKE ?", "%"+q.ServiceName+"%")
	}
	if q.Endpoint != "" {
		baseDB = baseDB.Where("endpoint ILIKE ?", "%"+q.Endpoint+"%")
	}
	if q.ReqHeader != "" {
		baseDB = baseDB.Where("req_header ILIKE ?", "%"+q.ReqHeader+"%")
	}
	if q.ReqMessage != "" {
		baseDB = baseDB.Where("req_message ILIKE ?", "%"+q.ReqMessage+"%")
	}
	if q.RespHeader != "" {
		baseDB = baseDB.Where("resp_header ILIKE ?", "%"+q.RespHeader+"%")
	}
	if q.RespMessage != "" {
		baseDB = baseDB.Where("resp_message ILIKE ?", "%"+q.RespMessage+"%")
	}
	if q.StatusCode != 0 {
		baseDB = baseDB.Where("status_code = ?", q.StatusCode)
	}
	if q.TimeUsage != 0 {
		baseDB = baseDB.Where("time_usage >= ?", q.TimeUsage)
	}

	if q.ReqDatetime != "" {
		if parsedDate, err := time.Parse("2006-01-02", q.ReqDatetime); err == nil {
			startOfDay := parsedDate
			endOfDay := parsedDate.Add(24 * time.Hour)
			baseDB = baseDB.Where("req_datetime >= ? AND req_datetime < ?", startOfDay, endOfDay)
		}
	}

	if q.RespDatetime != "" {
		if parsedDate, err := time.Parse("2006-01-02", q.RespDatetime); err == nil {
			startOfDay := parsedDate
			endOfDay := parsedDate.Add(24 * time.Hour)
			baseDB = baseDB.Where("resp_datetime >= ? AND resp_datetime < ?", startOfDay, endOfDay)
		}
	}

	var total int64
	if err := baseDB.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	sortField := q.SortingField
	if sortField == "" {
		sortField = "created_at"
	}
	sortOrder := strings.ToLower(q.SortingOrder)
	if sortOrder != "asc" {
		sortOrder = "desc"
	}

	if err := baseDB.Order(fmt.Sprintf("%s %s", sortField, sortOrder)).
		Limit(limit).
		Offset(offset).
		Find(&logs).Error; err != nil {
		return nil, 0, err
	}

	return logs, int(total), nil
}

func (repo *LogPg) SaveLoginEvent(ctx *gin.Context, accessLog *model.AccessLogs) error {
	if err := repo.db.WithContext(ctx).Create(accessLog).Error; err != nil {
		return err
	}
	return nil
}
