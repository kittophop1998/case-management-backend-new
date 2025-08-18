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

func (repo *UserPg) GetProfile(ctx *gin.Context, userId uuid.UUID) (*model.UserProfileResponse, error) {
	fmt.Println("Fetching user profile...")
	if userId == uuid.Nil {
		return nil, errors.New("unauthorized")
	}

	var user model.User
	if err := repo.db.WithContext(ctx).
		Preload("Role").
		Preload("Department").
		Preload("Section").
		Preload("Center").
		First(&user, model.User{ID: userId}).Error; err != nil {
		return nil, err
	}

	var perms []model.Permission
	err := repo.db.WithContext(ctx).
		Table("role_permissions AS rp").
		Select("p.id, p.key, p.name").
		Joins("JOIN permissions AS p ON rp.permission_id = p.id").
		Where("rp.role_id = ? AND rp.department_id = ?", user.RoleID, user.DepartmentID).
		Debug().Scan(&perms).Error
	if err != nil {
		return nil, err
	}

	resp := &model.UserProfileResponse{
		ID:       user.ID,
		Username: user.Username,
		Name:     user.Name,
		Role: model.Role{
			ID:   user.Role.ID,
			Name: user.Role.Name,
		},
		Department: model.Department{
			ID:   user.Department.ID,
			Name: user.Department.Name,
		},
		Section: model.Section{
			ID:   user.Section.ID,
			Name: user.Section.Name,
		},
		Center: model.Center{
			ID:   user.Center.ID,
			Name: user.Center.Name,
		},
		Permissions: perms,
	}

	fmt.Println("User profile retrieved:", resp)

	return resp, nil
}

func (repo *UserPg) GetAll(ctx *gin.Context, offset int, limit int, filter model.UserFilter) ([]*model.User, error) {
	var users []*model.User

	query := repo.db.WithContext(ctx).Model(&model.User{}).
		Preload("Role").
		Preload("Center").
		Preload("Section").
		Preload("Department").
		Joins("LEFT JOIN roles as role ON role.id = users.role_id").
		Joins("LEFT JOIN centers as center ON center.id = users.center_id").
		Joins("LEFT JOIN sections as section ON section.id = users.section_id").
		Joins("LEFT JOIN departments as department ON department.id = users.department_id")

	if err := query.Find(&users).Error; err != nil {
		return nil, err
	}

	if filter.Keyword != "" {
		kw := "%" + strings.TrimSpace(filter.Keyword) + "%"
		query = query.Where(
			repo.db.Where("users.name ILIKE ?", kw).
				Or("users.username ILIKE ?", kw),
		)
	}

	if filter.IsActive != nil {
		query = query.Where("users.is_active = ?", *filter.IsActive)
	}

	if filter.RoleID != uuid.Nil {
		query = query.Where("role.id = ?", filter.RoleID)
	}

	if filter.SectionID != uuid.Nil {
		query = query.Where("section.id = ?", filter.SectionID)
	}

	if filter.CenterID != uuid.Nil {
		query = query.Where("center.id = ?", filter.CenterID)
	}

	if filter.DepartmentID != uuid.Nil {
		query = query.Where("department.id = ?", filter.DepartmentID)
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
		Where(&id, model.User{ID: id}).Error; err != nil {
		return nil, err
	}

	return &user, nil
}

func (repo *UserPg) GetByUsername(ctx *gin.Context, username string) (*model.User, error) {
	var user model.User

	if err := repo.db.WithContext(ctx).
		Preload("Role").
		Preload("Center").
		Preload("Section").
		Where("username = ?", username).
		First(&user).Error; err != nil {
		return nil, err
	}

	return &user, nil
}

func (repo *UserPg) Create(ctx *gin.Context, user *model.CreateUpdateUserRequest) (uuid.UUID, error) {
	if err := repo.isDuplicate(ctx, "staff_id", *user.StaffID, uuid.Nil); err != nil {
		return uuid.Nil, err
	}

	if err := repo.isDuplicate(ctx, "username", user.Username, uuid.Nil); err != nil {
		return uuid.Nil, err
	}

	if err := repo.isDuplicate(ctx, "operator_id", *user.OperatorID, uuid.Nil); err != nil {
		return uuid.Nil, err
	}

	if err := repo.isDuplicate(ctx, "email", user.Email, uuid.Nil); err != nil {
		return uuid.Nil, err
	}

	isActive := true

	userToSave := &model.User{
		StaffID:      user.StaffID,
		Username:     user.Username,
		Name:         user.Name,
		OperatorID:   user.OperatorID,
		SectionID:    user.SectionID,
		CenterID:     user.CenterID,
		RoleID:       user.RoleID,
		DepartmentID: user.DepartmentID,
		Email:        user.Email,
		IsActive:     &isActive,
	}

	if err := repo.db.Create(&userToSave).Error; err != nil {
		if strings.Contains(err.Error(), "duplicate key value") {
			return uuid.Nil, fmt.Errorf("staffId %d already exists", user.StaffID)
		}

		return uuid.Nil, err
	}

	return userToSave.ID, nil
}

func (repo *UserPg) Update(ctx *gin.Context, id uuid.UUID, user model.CreateUpdateUserRequest) error {
	updateData := map[string]interface{}{}

	if user.Name != "" {
		updateData["name"] = user.Name
	}
	if user.RoleID != uuid.Nil {
		updateData["role_id"] = user.RoleID
	}
	if user.SectionID != uuid.Nil {
		updateData["section_id"] = user.SectionID
	}
	if user.CenterID != uuid.Nil {
		updateData["center_id"] = user.CenterID
	}
	if user.Username != "" {
		if err := repo.isDuplicate(ctx, "username", user.Username, id); err != nil {
			return err
		}
		updateData["username"] = user.Username
	}

	if user.DepartmentID != uuid.Nil {
		updateData["department_id"] = user.DepartmentID
	}

	if user.Email != "" {
		if err := repo.isDuplicate(ctx, "email", user.Email, id); err != nil {
			return err
		}
		updateData["email"] = user.Email
	}

	if user.StaffID != nil {
		if err := repo.isDuplicate(ctx, "staff_id", *user.StaffID, id); err != nil {
			return err
		}
		updateData["staff_id"] = *user.StaffID
	}

	if user.OperatorID != nil {
		if err := repo.isDuplicate(ctx, "operator_id", *user.OperatorID, id); err != nil {
			return err
		}
		updateData["operator_id"] = *user.OperatorID
	}

	if user.IsActive != nil {
		updateData["is_active"] = *user.IsActive
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

func (repo *UserPg) isDuplicate(ctx *gin.Context, field string, value interface{}, excludeID uuid.UUID) error {
	var existingUser model.User
	query := repo.db.WithContext(ctx).Where(fmt.Sprintf("%s = ?", field), value)

	if excludeID != uuid.Nil {
		query = query.Where("id <> ?", excludeID)
	}

	if err := query.First(&existingUser).Error; err == nil {
		return fmt.Errorf("%s %v already exists", field, value)
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}

	return nil
}
