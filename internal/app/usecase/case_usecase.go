package usecase

import (
	"case-management/internal/domain/model"
	"case-management/internal/domain/repository"
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type CaseUseCase struct {
	repo repository.CaseRepository
}

func NewCaseUseCase(repo repository.CaseRepository) *CaseUseCase {
	return &CaseUseCase{repo: repo}
}

func (uc *CaseUseCase) GetAllCases(ctx *gin.Context, page, limit int, category string, currID uuid.UUID) ([]*model.CaseResponse, int, error) {
	offset := (page - 1) * limit
	caseRepo, total, err := uc.repo.GetAllCase(ctx, offset, limit, category, currID)
	if err != nil {
		return nil, 0, err
	}

	var caseResponses []*model.CaseResponse
	for _, c := range caseRepo {

		caseResponses = append(caseResponses, &model.CaseResponse{
			CustomerID:   c.CustomerID,
			CustomerName: c.CustomerName,
			Status:       c.Status.Name,
			CaseType:     c.CaseType.Name,
			CurrentQueue: c.Queue.Name,
			CurrentUser:  fmt.Sprintf("%s - %s", c.AssignedToUser.Name, c.AssignedToUser.Center.Name),
			SLADate:      c.EndDate.String(),
			CreateDate:   c.CreatedAt.String(),
			CaseID:       c.ID.String(),
			AeonID:       c.AeonID,
			CaseGroup:    "General",
			CreatedBy:    c.CreatedBy.String(),
			CreatedDate:  c.CreatedAt.String(),
			CasePriority: c.Priority,
			ClosedDate:   c.ClosedDate.String(),
			ReceivedFrom: "Fraud",
		})
	}
	return caseResponses, total, nil

}

func (uc *CaseUseCase) GetCaseByID(c *gin.Context, id string) (*model.Cases, error) {
	caseID, err := uuid.Parse(id)
	if err != nil {
		return nil, err
	}
	return uc.repo.GetCaseByID(c, caseID)
}

func (uc *CaseUseCase) CreateCaseInquiry(ctx *gin.Context, createdByID uuid.UUID, caseReq *model.CreateCaseInquiryRequest) (uuid.UUID, error) {
	statusMap, _ := uc.repo.LoadCaseStatus(ctx)

	caseToSave := &model.Cases{
		CaseTypeID:        caseReq.CaseTypeID,
		CustomerName:      caseReq.CustomerName,
		CustomerID:        caseReq.CustomerID,
		DispositionMainID: caseReq.DispositionMainID,
		DispositionSubID:  caseReq.DispositionSubID,
		AssignedToUserID:  createdByID,
		ProductID:         caseReq.ProductID,
		Priority:          "Normal",
		StatusID:          statusMap["created"],
		StartDate:         time.Now(),
		EndDate:           time.Now().Add(72 * time.Hour), // SLA 3 days
		CreatedBy:         createdByID,
		UpdatedBy:         createdByID,
	}

	caseId, err := uc.repo.CreateCaseInquiry(ctx, caseToSave)
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

func (uc *CaseUseCase) AddInitialDescription(c *gin.Context, caseID string, newDescription string) error {
	caseUUID, err := uuid.Parse(caseID)
	if err != nil {
		return err
	}
	return uc.repo.AddInitialDescription(c, caseUUID, newDescription)
}

func (uc *CaseUseCase) GetAllDisposition(ctx *gin.Context) ([]model.DispositionItem, error) {
	return uc.repo.GetAllDisposition(ctx)
}

func (uc *CaseUseCase) AddCaseNote(ctx *gin.Context, createdByID uuid.UUID, caseID uuid.UUID, input *model.CaseNoteRequest) (uuid.UUID, error) {
	note := &model.CaseNotes{
		CaseId:  caseID,
		UserId:  createdByID,
		Content: input.Content,
	}

	return uc.repo.AddCaseNote(ctx, note)
}
