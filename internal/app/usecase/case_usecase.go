package usecase

import (
	"case-management/internal/domain/model"
	"case-management/internal/domain/repository"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type CaseUseCase struct {
	repo repository.CaseRepository
}

func NewCaseUseCase(repo repository.CaseRepository) *CaseUseCase {
	return &CaseUseCase{repo: repo}
}

func (uc *CaseUseCase) GetAllCases(c *gin.Context) ([]*model.Cases, error) {
	return uc.repo.GetAllCase(c, 10, 0, model.CaseFilter{})
}

func (uc *CaseUseCase) GetCaseByID(c *gin.Context, id string) (*model.Cases, error) {
	caseID, err := uuid.Parse(id)
	if err != nil {
		return nil, err
	}
	return uc.repo.GetCaseByID(c, caseID)
}

func (uc *CaseUseCase) CreateCase(c *gin.Context, caseData *model.CreateCaseRequest) (uuid.UUID, error) {
	// Create Case
	caseId, err := uc.repo.CreateCase(c, caseData)
	if err != nil {
		return uuid.Nil, err
	}

	// Save Case Disposition Main
	// err := uc.repo.CreateCaseDispositionMains(c, caseData.DispositionMains)
	// if err != nil {
	// 	return uuid.Nil, err
	// }

	// // Save Case Disposition Subs
	// err = uc.repo.CreateCaseDispositionSubs(c, caseData.DispositionSubs)
	// if err != nil {
	// 	return uuid.Nil, err
	// }

	return caseId, nil
}

// func (uc *CaseUseCase) CreateNoteType(c *gin.Context, note model.NoteTypes) (*model.NoteTypes, error) {
// 	return uc.repo.CreateNoteType(c, note)
// }

// func (uc *CaseUseCase) GetNoteTypeByID(c *gin.Context, noteTypeID string) (*model.NoteTypes, error) {
// 	noteTypeUUID, err := uuid.Parse(noteTypeID)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return uc.repo.GetNoteTypeByID(c, noteTypeUUID)
// }

func (uc *CaseUseCase) AddInitialDescription(c *gin.Context, caseID string, newDescription string) error {
	caseUUID, err := uuid.Parse(caseID)
	if err != nil {
		return err
	}
	return uc.repo.AddInitialDescription(c, caseUUID, newDescription)
}

func (uc *CaseUseCase) GetAllDisposition(ctx *gin.Context, page, limit int, keyword string) ([]model.DispositionMain, int, error) {
	offset := (page - 1) * limit
	return uc.repo.GetAllDisposition(ctx, limit, offset, keyword)
}
