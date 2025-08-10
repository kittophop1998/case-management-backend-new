package database

import (
	"case-management/internal/domain/model"

	"gorm.io/gorm"
)

type CustomerPg struct {
	db *gorm.DB
}

func NewCustomerPg(db *gorm.DB) *CustomerPg {
	return &CustomerPg{db: db}
}

func (c *CustomerPg) CreateCustomerNote(note *model.CustomerNote) error {
	if err := c.db.Create(note).Error; err != nil {
		return err
	}
	return nil
}

func (c *CustomerPg) GetAllCustomerNotes(customerID string) ([]*model.CustomerNote, error) {
	var notes []*model.CustomerNote
	if err := c.db.Where("customer_id = ?", customerID).Find(&notes).Error; err != nil {
		return nil, err
	}
	return notes, nil
}
