package database

import (
	"case-management/internal/domain/model"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

const tableUsers = "users"

type UserPg struct {
	db *gorm.DB
}

func NewUserPg(db *gorm.DB) *UserPg {
	return &UserPg{db: db}
}

func (u *UserPg) FindAll(c *gin.Context) ([]model.User, error) {
	var users []model.User
	if err := u.db.WithContext(c).Table(tableUsers).Find(&users).Error; err != nil {
		return nil, err
	}
	return users, nil
}
