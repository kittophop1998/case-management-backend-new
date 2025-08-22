package usecase

import (
	"case-management/internal/domain/model"
	"case-management/internal/domain/repository"
	"time"

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

func (uc *CustomerUseCase) GetAllCustomerNotes(ctx *gin.Context, customerID string, page, limit int) ([]model.CustomerNoteResponse, int, error) {
	loc := time.FixedZone("Asia/Bangkok", 7*60*60) // +7 ชั่วโมง
	offset := (page - 1) * limit

	notes, total, err := uc.CustomerRepo.GetAllCustomerNotes(ctx, customerID, limit, offset)
	if err != nil {
		return nil, 0, err
	}

	results := make([]model.CustomerNoteResponse, len(notes))
	for i, note := range notes {
		results[i] = model.CustomerNoteResponse{
			ID:          note.ID.String(),
			NoteType:    note.NoteType.Name,
			NoteDetail:  note.Note,
			CreatedBy:   note.CreatedBy,
			CreatedDate: note.CreatedAt.In(loc).Format("02 Jan 2006 15:04:05"),
		}
	}

	return results, total, nil
}

func (uc *CustomerUseCase) GetNoteTypes(ctx *gin.Context) ([]*model.NoteTypes, error) {
	return uc.CustomerRepo.GetNoteTypes(ctx)
}

func (uc *CustomerUseCase) CountNotes(ctx *gin.Context, customerID string) (int, error) {
	return uc.CustomerRepo.CountNotes(ctx, customerID)
}
