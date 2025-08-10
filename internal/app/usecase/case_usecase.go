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
	return uc.repo.CreateCase(c, caseData)
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

// func (uc *CaseUseCase) CountCasesWithFilter(c *gin.Context, filter model.CaseFilter) (int, error) {
// 	return uc.repo.CountWithFilter(c, filter)
// }
