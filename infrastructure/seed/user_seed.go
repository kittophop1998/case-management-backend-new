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
	defaultPassword := "Aeon*123"
	userType := "local"
	staffID := uint(1)
	operatorID := uint(1)
	users := []model.User{
		{
			Name:         "Admin",
			Username:     "admin",
			Password:     defaultPassword,
			UserTypes:    userType,
			SectionID:    sectionMap["Inbound"],
			CenterID:     centerMap["BKK"],
			RoleID:       roleMap["Admin"],
			StaffID:      &staffID,
			IsActive:     &isActive,
			Email:        "admin@admin.com",
			OperatorID:   &operatorID,
			DepartmentID: departmentMap["Marketing"],
		},
		{
			Name:         "System-support",
			Username:     "support",
			Password:     defaultPassword,
			UserTypes:    userType,
			SectionID:    sectionMap["Inbound"],
			CenterID:     centerMap["BKK"],
			RoleID:       roleMap["Admin"],
			StaffID:      &staffID,
			IsActive:     &isActive,
			Email:        "admin@admin.com",
			OperatorID:   &operatorID,
			DepartmentID: departmentMap["Marketing"],
		},
	}

	for _, user := range users {
		var existingUser model.User
		if err := db.Where("username = ?", user.Username).First(&existingUser).Error; err != nil {
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
