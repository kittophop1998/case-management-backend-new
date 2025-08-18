package database

import (
	"case-management/internal/domain/model"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type MasterDataPg struct {
	db *gorm.DB
}

func NewMasterDataPg(db *gorm.DB) *MasterDataPg {
	return &MasterDataPg{db: db}
}

func (repo *MasterDataPg) FindAll(ctx *gin.Context) (map[string]interface{}, error) {
	result := make(map[string]interface{})

	var sections []model.Section
	if err := repo.db.WithContext(ctx).Find(&sections).Error; err != nil {
		return nil, err
	}
	result["sections"] = sections

	var centers []model.Center
	if err := repo.db.WithContext(ctx).Find(&centers).Error; err != nil {
		return nil, err
	}
	result["centers"] = centers

	var permissions []model.Permission
	if err := repo.db.WithContext(ctx).Find(&permissions).Error; err != nil {
		return nil, err
	}
	result["permissions"] = permissions

	var roles []model.Role
	if err := repo.db.WithContext(ctx).Find(&roles).Error; err != nil {
		return nil, err
	}
	result["roles"] = roles

	var departments []model.Department
	if err := repo.db.WithContext(ctx).Find(&departments).Error; err != nil {
		return nil, err
	}
	result["departments"] = departments

	return result, nil
}
