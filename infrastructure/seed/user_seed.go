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
	staffIdAdmin := uint(1)
	operatorIdAdmin := uint(1)
	// staffIdSupport := uint(2)
	// operatorIdSupport := uint(2)
	users := []model.User{
		{
			Name:         "Admin",
			Username:     "admin",
			Password:     defaultPassword,
			UserTypes:    userType,
			SectionID:    sectionMap["CHL"],
			CenterID:     centerMap["BKK"],
			RoleID:       roleMap["Admin"],
			StaffID:      &staffIdAdmin,
			IsActive:     &isActive,
			Email:        "admin@admin.com",
			OperatorID:   &operatorIdAdmin,
			DepartmentID: departmentMap["Marketing"],
		},
		// {
		// 	Name:         "System-support",
		// 	Username:     "support",
		// 	Password:     defaultPassword,
		// 	UserTypes:    userType,
		// 	SectionID:    sectionMap["CHL"],
		// 	CenterID:     centerMap["BKK"],
		// 	RoleID:       roleMap["Admin"],
		// 	StaffID:      &staffIdSupport,
		// 	IsActive:     &isActive,
		// 	Email:        "support@admin.com",
		// 	OperatorID:   &operatorIdSupport,
		// 	DepartmentID: departmentMap["Marketing"],
		// },
	}

	for _, user := range users {
		if err := db.FirstOrCreate(&model.User{}, user).Error; err != nil {
			panic("Failed to seed user: " + err.Error())
		}
	}
}
