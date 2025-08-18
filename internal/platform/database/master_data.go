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
	if err := repo.db.WithContext(ctx).Model(&model.Section{}).Find(&sections).Error; err != nil {
		return nil, err
	}
	result["sections"] = sections

	// var roles []model.Role
	// if err := repo.db.WithContext(ctx).Preload("Permissions").Find(&roles).Error; err != nil {
	// 	return nil, err
	// }
	// result["roles"] = roles

	var centers []model.Center
	if err := repo.db.WithContext(ctx).Model(&model.Center{}).Find(&centers).Error; err != nil {
		return nil, err
	}
	result["centers"] = centers

	// var permissions []model.Permission
	// if err := repo.db.WithContext(ctx).Model(&model.Permission{}).Find(&permissions).Error; err != nil {
	// 	return nil, err
	// }
	// result["permissions"] = permissions

	var departments []model.Department
	if err := repo.db.WithContext(ctx).Model(&model.Department{}).Find(&departments).Error; err != nil {
		return nil, err
	}
	result["departments"] = departments

	return result, nil
}
