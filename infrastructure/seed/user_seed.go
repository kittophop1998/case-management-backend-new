package seed

import (
	"case-management/internal/domain/model"
	"fmt"

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
	staffIdAdmin := uint(1)
	operatorIdAdmin := uint(1)
	staffIdSupport := uint(2)
	operatorIdSupport := uint(2)
	users := []model.User{
		{
			Name:         "Admin",
			Username:     "admin",
			Password:     defaultPassword,
			UserTypes:    userType,
			SectionID:    sectionMap["System"],
			CenterID:     centerMap["BKK"],
			RoleID:       roleMap["System"],
			StaffID:      &staffIdAdmin,
			IsActive:     &isActive,
			Email:        "admin@admin.com",
			OperatorID:   &operatorIdAdmin,
			DepartmentID: departmentMap["System"],
		},
		{
			Name:         "System-support",
			Username:     "support",
			Password:     defaultPassword,
			UserTypes:    userType,
			SectionID:    sectionMap["System"],
			CenterID:     centerMap["BKK"],
			RoleID:       roleMap["System"],
			StaffID:      &staffIdSupport,
			IsActive:     &isActive,
			Email:        "support@admin.com",
			OperatorID:   &operatorIdSupport,
			DepartmentID: departmentMap["System"],
		},
	}

	for i := 3; i < 22; i++ {
		staffId := uint(i)
		operatorId := uint(i)
		users = append(users, model.User{
			Name:         fmt.Sprintf("User-%d", i),
			Username:     fmt.Sprintf("user%d", i),
			Password:     defaultPassword,
			UserTypes:    userType,
			SectionID:    sectionMap["System"],
			CenterID:     centerMap["BKK"],
			RoleID:       roleMap["Admin"],
			StaffID:      &staffId,
			IsActive:     &isActive,
			Email:        fmt.Sprintf("user%d@admin.com", i),
			OperatorID:   &operatorId,
			DepartmentID: departmentMap["Finance"],
		})
	}

	for _, user := range users {
		if err := db.FirstOrCreate(&model.User{}, user).Error; err != nil {
			panic("Failed to seed user: " + err.Error())
		}
	}
}
