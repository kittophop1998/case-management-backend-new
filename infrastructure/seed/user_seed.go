package seed

import (
	"case-management/internal/domain/model"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

func SeedUser(
	db *gorm.DB,
	roleMap map[string]uuid.UUID,
	sectionMap map[string]uuid.UUID,
	centerMap map[string]uuid.UUID,
	departmentMap map[string]uuid.UUID,
	queueMap map[string]uuid.UUID,
) {
	isActive := true
	users := []model.User{
		{
			Name:         "admin",
			SectionID:    sectionMap["Inbound"],
			CenterID:     centerMap["BKK"],
			RoleID:       roleMap["Admin"],
			QueueID:      queueMap["Inbound Agent Supervisor BKK Queue"],
			AgentID:      1,
			IsActive:     &isActive,
			Email:        "admin@admin.com",
			OperatorID:   1,
			DepartmentID: departmentMap["Marketing"],
		},
	}

	for _, user := range users {
		var existingUser model.User
		if err := db.Where("name = ?", user.Name).First(&existingUser).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				// User does not exist, create it
				if err := db.Create(&user).Error; err != nil {
					panic("Failed to seed user: " + err.Error())
				}
			} else {
				panic("Failed to check existing user: " + err.Error())
			}
		}
	}
}
