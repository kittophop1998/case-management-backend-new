package repository

import "case-management/internal/domain/model"

type CustomerRepository interface {
	CreateCustomerNote(note *model.CustomerNote) error
	GetAllCustomerNotes(customerID string) ([]*model.CustomerNote, error)
}
