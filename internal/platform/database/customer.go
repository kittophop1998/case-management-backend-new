package database

import (
	"case-management/internal/domain/model"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type CustomerPg struct {
	db *gorm.DB
}

func NewCustomerPg(db *gorm.DB) *CustomerPg {
	return &CustomerPg{db: db}
}

func (c *CustomerPg) CreateCustomerNote(ctx *gin.Context, note *model.CustomerNote) error {
	if err := c.db.WithContext(ctx).Create(note).Error; err != nil {
		return err
	}
	return nil
}

func (c *CustomerPg) GetAllCustomerNotes(ctx *gin.Context, customerID string, limit, offset int) ([]*model.CustomerNote, int, error) {
	var notes []*model.CustomerNote
	if err := c.db.WithContext(ctx).
		Preload("NoteType").
		Where("customer_id = ?", customerID).
		Limit(limit).
		Offset(offset).
		Find(&notes).Error; err != nil {
		return nil, 0, err
	}

	var total int64
	if err := c.db.WithContext(ctx).Model(&model.CustomerNote{}).Where("customer_id = ?", customerID).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	return notes, int(total), nil
}

func (c *CustomerPg) GetNoteTypes(ctx *gin.Context) ([]*model.NoteTypes, error) {
	var noteTypes []*model.NoteTypes
	if err := c.db.WithContext(ctx).Find(&noteTypes).Error; err != nil {
		return nil, err
	}
	return noteTypes, nil
}

func (c *CustomerPg) CountNotes(ctx *gin.Context, customerID string) (int, error) {
	var count int64
	if err := c.db.WithContext(ctx).Model(&model.CustomerNote{}).Where("customer_id = ?", customerID).Debug().Count(&count).Error; err != nil {
		return 0, err
	}
	return int(count), nil
}
