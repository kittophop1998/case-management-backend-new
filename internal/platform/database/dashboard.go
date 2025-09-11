package database

import (
	"case-management/internal/domain/model"
	"context"
	"log"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type dashboardPg struct {
	db *gorm.DB
}

func NewDashboardPg(db *gorm.DB) *dashboardPg {
	return &dashboardPg{db: db}
}

func (r *dashboardPg) SaveApiLog(ctx context.Context, logData *model.ApiLogs) error {
	log.Println("SaveApiLog called in repo")
	log.Println("DB connection is nil?", r.db == nil)

	logData.ID = uuid.New().String()
	logData.CreatedAt = time.Now()

	err := r.db.Create(logData).Error
	if err != nil {
		log.Println("DB insert error:", err)
	}
	return err
}
