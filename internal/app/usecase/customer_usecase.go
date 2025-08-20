package usecase

import (
	"case-management/internal/domain/model"
	"case-management/internal/domain/repository"

	"github.com/gin-gonic/gin"
)

type CustomerUseCase struct {
	CustomerRepo repository.CustomerRepository
}

func NewCustomerUseCase(repo repository.CustomerRepository) *CustomerUseCase {
	return &CustomerUseCase{CustomerRepo: repo}
}

func (uc *CustomerUseCase) CreateCustomerNote(ctx *gin.Context, note *model.CustomerNote) error {
	return uc.CustomerRepo.CreateCustomerNote(ctx, note)
}

func (uc *CustomerUseCase) GetAllCustomerNotes(ctx *gin.Context, customerID string) ([]*model.CustomerNote, error) {
	return uc.CustomerRepo.GetAllCustomerNotes(ctx, customerID)
}

func (uc *CustomerUseCase) GetNoteTypes(ctx *gin.Context) ([]*model.NoteTypes, error) {
	return uc.CustomerRepo.GetNoteTypes(ctx)
}

func (uc *CustomerUseCase) CountNotes(ctx *gin.Context, customerID string) (int, error) {
	return uc.CustomerRepo.CountNotes(ctx, customerID)
}
