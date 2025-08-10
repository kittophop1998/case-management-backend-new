package usecase

import (
	"case-management/internal/domain/model"
	"case-management/internal/domain/repository"
)

type CustomerUseCase struct {
	CustomerRepo repository.CustomerRepository
}

func NewCustomerUseCase(repo repository.CustomerRepository) *CustomerUseCase {
	return &CustomerUseCase{CustomerRepo: repo}
}

func (uc *CustomerUseCase) CreateCustomerNote(note *model.CustomerNote) error {
	return uc.CustomerRepo.CreateCustomerNote(note)
}

func (uc *CustomerUseCase) GetAllCustomerNotes(customerID string) ([]*model.CustomerNote, error) {
	return uc.CustomerRepo.GetAllCustomerNotes(customerID)
}
