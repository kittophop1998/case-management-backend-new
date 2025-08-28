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
) {
	isActive := true
	defaultPassword := "aeon"
	userType := "local"
	staffId1 := uint(1)
	operatorId1 := uint(1)
	staffId2 := uint(2)
	operatorId2 := uint(2)
	staffId3 := uint(3)
	operatorId3 := uint(3)
	staffId4 := uint(4)
	operatorId4 := uint(4)
	staffId5 := uint(5)
	operatorId5 := uint(5)
	staffId6 := uint(6)
	operatorId6 := uint(6)

	users := []model.User{
		{
			Name:         "Admin",
			Username:     "admin",
			Password:     defaultPassword,
			UserTypes:    userType,
			SectionID:    sectionMap["SYSTEM"],
			CenterID:     centerMap["BKK"],
			RoleID:       roleMap["System"],
			StaffID:      &staffId1,
			IsActive:     &isActive,
			Email:        "admin@admin.com",
			OperatorID:   &operatorId1,
			DepartmentID: departmentMap["SYSTEM"],
		},
		{
			Name:         "System-Support",
			Username:     "support",
			Password:     defaultPassword,
			UserTypes:    userType,
			SectionID:    sectionMap["SYSTEM"],
			CenterID:     centerMap["BKK"],
			RoleID:       roleMap["System"],
			StaffID:      &staffId2,
			IsActive:     &isActive,
			Email:        "support@admin.com",
			OperatorID:   &operatorId2,
			DepartmentID: departmentMap["SYSTEM"],
		},
		{
			Name:         "สมชาย ใจดี",
			Username:     "somchai",
			Password:     defaultPassword,
			UserTypes:    userType,
			SectionID:    sectionMap["ONB"],
			CenterID:     centerMap["BKK"],
			RoleID:       roleMap["Admin"],
			StaffID:      &staffId3,
			IsActive:     &isActive,
			Email:        "somchai@admin.com",
			OperatorID:   &operatorId3,
			DepartmentID: departmentMap["CMS"],
		},
		{
			Name:         "สมหญิง สายสมร",
			Username:     "somying",
			Password:     defaultPassword,
			UserTypes:    userType,
			SectionID:    sectionMap["ONB"],
			CenterID:     centerMap["BKK"],
			RoleID:       roleMap["Staff"],
			StaffID:      &staffId4,
			IsActive:     &isActive,
			Email:        "somying@admin.com",
			OperatorID:   &operatorId4,
			DepartmentID: departmentMap["CMS"],
		},
		{
			Name:         "วิชัย พึ่งพา",
			Username:     "wichai",
			Password:     defaultPassword,
			UserTypes:    userType,
			SectionID:    sectionMap["ONB"],
			CenterID:     centerMap["BKK"],
			RoleID:       roleMap["Supervisor"],
			StaffID:      &staffId5,
			IsActive:     &isActive,
			Email:        "wichai@admin.com",
			OperatorID:   &operatorId5,
			DepartmentID: departmentMap["CMS"],
		},
		{
			Name:         "กมลทิพย์ ศรีสุข",
			Username:     "kamoltip",
			Password:     defaultPassword,
			UserTypes:    userType,
			SectionID:    sectionMap["ONB"],
			CenterID:     centerMap["BKK"],
			RoleID:       roleMap["AsstManager Up"],
			StaffID:      &staffId6,
			IsActive:     &isActive,
			Email:        "kamoltip@admin.com",
			OperatorID:   &operatorId6,
			DepartmentID: departmentMap["CMS"],
		},
	}

	for i := 7; i < 30; i++ {
		staffId := uint(i)
		operatorId := uint(i)
		users = append(users, model.User{
			Name:         fmt.Sprintf("User-%d", i),
			Username:     fmt.Sprintf("user%d", i),
			Password:     defaultPassword,
			UserTypes:    userType,
			SectionID:    sectionMap["ONB"],
			CenterID:     centerMap["BKK"],
			RoleID:       roleMap["Admin"],
			StaffID:      &staffId,
			IsActive:     &isActive,
			Email:        fmt.Sprintf("user%d@admin.com", i),
			OperatorID:   &operatorId,
			DepartmentID: departmentMap["CMS"],
		})
	}

	for _, user := range users {
		if err := db.FirstOrCreate(&user, model.User{Username: user.Username, StaffID: user.StaffID}).Error; err != nil {
			panic("Failed to seed user: " + err.Error())
		}
	}
}
