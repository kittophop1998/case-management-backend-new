package seed

import (
	"case-management/internal/domain/model"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

func SeedReasonCodes(db *gorm.DB) {
	reasonCodes := []model.ReasonCode{
		{
			Code:              "RC001",
			DescriptionEn:     "Customer Request",
			DescriptionTh:     "คำขอของลูกค้า",
			Category:          "General",
			SLAResponseTime:   "1 hour",
			SLAResolutionTime: "24 hours",
			Note:              "Initial customer request",
			CreatedBy:         uuid.New(),
			CreatedAt:         time.Now(),
		},
	}

	if err := db.Create(&reasonCodes).Error; err != nil {
		panic("Failed to seed reason codes: " + err.Error())
	}
}
