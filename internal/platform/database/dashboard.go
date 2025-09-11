package database

import (
	"case-management/internal/domain/model"
	"context"
	"fmt"
	"log"
	"strings"
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

func (r *dashboardPg) GetCustIDByAeonID(ctx context.Context, inputID string) (string, error) {

	if isCustID(inputID) {
		return inputID, nil
	}

	var custID string

	escapedAeonID := strings.ReplaceAll(inputID, "'", "''")

	// dblink query
	dblinkQuery := fmt.Sprintf(
		"SELECT cust_id, aeon_id FROM d_mobile_app_daily WHERE aeon_id = '%s' LIMIT 1",
		escapedAeonID,
	)
	escapedDblinkQuery := strings.ReplaceAll(dblinkQuery, "'", "''")

	query := fmt.Sprintf(`
        SELECT cust_id FROM
        dblink(
            'host=10.251.1.103 dbname=idmapping user=bff password=Aeon*123',
            '%s'
        ) AS t(cust_id varchar, aeon_id varchar)
        WHERE aeon_id = '%s';
    `, escapedDblinkQuery, escapedAeonID)

	fmt.Println("dblink query:\n", query)

	if err := r.db.WithContext(ctx).Raw(query).Scan(&custID).Error; err != nil {
		return "", fmt.Errorf("failed to get cust_id from aeon_id via dblink: %w", err)
	}

	if custID == "" {
		return "", fmt.Errorf("cust_id not found for aeon_id: %s", inputID)
	}

	return custID, nil
}

func isCustID(id string) bool {
	if len(id) != 13 {
		return false
	}
	for _, r := range id {
		if r < '0' || r > '9' {
			return false
		}
	}
	return true
}
