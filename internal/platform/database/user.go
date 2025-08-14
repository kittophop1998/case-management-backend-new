package database

import (
	"case-management/internal/domain/model"
	"errors"
	"fmt"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserPg struct {
	db *gorm.DB
}

func NewUserPg(db *gorm.DB) *UserPg {
	return &UserPg{db: db}
}

func (repo *UserPg) GetAll(ctx *gin.Context, offset int, limit int, filter model.UserFilter) ([]*model.User, error) {
	var users []*model.User

	query := repo.db.WithContext(ctx).Model(&model.User{}).
		Preload("Role").
		Preload("Center").
		Preload("Section").
		Preload("Department").
		Joins("LEFT JOIN roles ON roles.id = users.role_id").
		Joins("LEFT JOIN centers ON centers.id = users.center_id").
		Joins("LEFT JOIN sections ON sections.id = users.section_id").
		Joins("LEFT JOIN departments ON departments.id = users.department_id")

	if err := query.Find(&users).Error; err != nil {
		return nil, err
	}

	if filter.Keyword != "" {
		kw := "%" + strings.TrimSpace(filter.Keyword) + "%"
		query = query.Where(
			repo.db.Where("users.name ILIKE ?", kw).
				Or("users.username ILIKE ?", kw).
				Or("users.email ILIKE ?", kw).
				Or("CAST(users.agent_id AS TEXT) ILIKE ?", kw),
		)
	}

	if filter.IsActive != nil {
		query = query.Where("users.is_active = ?", *filter.IsActive)
	}

	if filter.RoleID != uuid.Nil {
		query = query.Where("roles.id = ?", filter.RoleID)
	}

	if filter.SectionID != uuid.Nil {
		query = query.Where("sections.id = ?", filter.SectionID)
	}

	if filter.CenterID != uuid.Nil {
		query = query.Where("centers.id = ?", filter.CenterID)
	}

	if filter.Sort != "" {
		query = query.Order(filter.Sort)
	}

	if err := query.Limit(limit).Offset(offset).Find(&users).Error; err != nil {
		return nil, err
	}

	return users, nil
}

func (repo *UserPg) GetById(ctx *gin.Context, id uuid.UUID) (*model.User, error) {
	var user model.User

	if err := repo.db.WithContext(ctx).
		Preload("Role").
		Preload("Center").
		Preload("Section").
		Preload("Department").
		Preload("Role.Permissions").
		Where("id = ?", id).First(&user).Error; err != nil {
		return nil, err
	}

	return &user, nil
}

func (repo *UserPg) GetByUsername(ctx *gin.Context, username string) (*model.User, error) {
	var user model.User

	if err := repo.db.WithContext(ctx).
		Preload("Role").
		Preload("Center").
		Preload("Role.Permissions").
		Preload("Section").
		Where("username = ?", username).
		First(&user).Error; err != nil {
		return nil, err
	}

	return &user, nil
}

func (repo *UserPg) Create(ctx *gin.Context, user *model.CreateUpdateUserRequest) (uuid.UUID, error) {
	var existingUser model.User

	if err := repo.db.WithContext(ctx).Where("staff_id = ?", user.StaffID).First(&existingUser).Error; err == nil {
		return uuid.Nil, fmt.Errorf("staffId %d already exists", user.StaffID)
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		return uuid.Nil, err
	}

	if err := repo.db.WithContext(ctx).Where("operator_id = ?", user.OperatorID).First(&existingUser).Error; err == nil {
		return uuid.Nil, fmt.Errorf("operatorId %d already exists", user.OperatorID)
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		return uuid.Nil, err
	}

	isActive := true

	userToSave := &model.User{
		StaffID:      user.StaffID,
		Name:         user.Name,
		OperatorID:   user.OperatorID,
		SectionID:    user.SectionID,
		CenterID:     user.CenterID,
		RoleID:       user.RoleID,
		DepartmentID: user.DepartmentID,
		Email:        user.Email,
		IsActive:     &isActive,
	}

	parts := strings.Split(user.Email, "@")
	if len(parts) == 2 {
		userToSave.Username = parts[0]
	} else {
		userToSave.Username = ""
	}

	if err := repo.db.Create(&userToSave).Error; err != nil {
		if strings.Contains(err.Error(), "duplicate key value") {
			return uuid.Nil, fmt.Errorf("staffId %d already exists", user.StaffID)
		}

		return uuid.Nil, err
	}

	return userToSave.ID, nil
}

func (repo *UserPg) Update(ctx *gin.Context, id uuid.UUID, input model.CreateUpdateUserRequest) error {
	updateData := map[string]interface{}{}

	if input.Name != "" {
		updateData["name"] = input.Name
	}
	if input.RoleID != uuid.Nil {
		updateData["role_id"] = input.RoleID
	}
	if input.SectionID != uuid.Nil {
		updateData["section_id"] = input.SectionID
	}
	if input.CenterID != uuid.Nil {
		updateData["center_id"] = input.CenterID
	}
	if input.Username != "" {
		updateData["username"] = input.Username
	}

	if input.DepartmentID != uuid.Nil {
		updateData["department_id"] = input.DepartmentID
	}

	if input.Email != "" {
		updateData["email"] = input.Email
	}

	if input.StaffID != nil {
		updateData["staff_id"] = *input.StaffID
	}
	if input.OperatorID != nil {
		updateData["operator_id"] = *input.OperatorID
	}

	if input.IsActive != nil {
		updateData["is_active"] = *input.IsActive
	}

	if len(updateData) == 0 {
		return errors.New("no valid fields to update")
	}

	if err := repo.db.WithContext(ctx).Model(&model.User{}).Where("id = ?", id).Debug().Updates(updateData).Error; err != nil {
		return err
	}

	return nil
}

func (repo *UserPg) Count(ctx *gin.Context) (int, error) {
	var count int64
	if err := repo.db.WithContext(ctx).Model(&model.User{}).Count(&count).Error; err != nil {
		return 0, err
	}
	return int(count), nil
}

func (repo *UserPg) CountWithFilter(ctx *gin.Context, filter model.UserFilter) (int, error) {
	var count int64
	query := repo.db.WithContext(ctx).Model(&model.User{}).
		Joins("LEFT JOIN roles ON roles.id = users.role_id").
		Joins("LEFT JOIN centers ON centers.id = users.center_id").
		Joins("LEFT JOIN sections ON sections.id = users.section_id")

	if filter.IsActive != nil {
		query = query.Where("users.is_active = ?", *filter.IsActive)
	}

	if filter.Role != "" {
		query = query.Where("roles.name = ?", filter.Role)
	}

	if filter.RoleID != uuid.Nil {
		query = query.Where("roles.id = ?", filter.RoleID)
	}

	if filter.Section.Name != "" {
		query = query.Where("sections.name = ?", strings.TrimSpace(filter.Section.Name))
	}

	if filter.SectionID != uuid.Nil {
		query = query.Where("sections.id = ?", filter.SectionID)
	}

	if filter.Center != "" {
		query = query.Where("centers.name ILIKE ?", "%"+strings.TrimSpace(filter.Center)+"%")
	}

	if filter.CenterID != uuid.Nil {
		query = query.Where("centers.id = ?", filter.CenterID)
	}

	if err := query.Count(&count).Error; err != nil {
		return 0, err
	}

	return int(count), nil
}
