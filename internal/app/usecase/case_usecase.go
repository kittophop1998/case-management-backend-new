package usecase

import (
	"case-management/internal/domain/model"
	"case-management/internal/domain/repository"
	"case-management/utils"
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

func (uc *CaseUseCase) GetCaseByID(c *gin.Context, id uuid.UUID) (*model.CaseDetailResponse, error) {
	caseData, err := uc.repo.GetCaseByID(c, id)
	if err != nil {
		return nil, err
	}

	caseDetail := &model.CaseDetailResponse{
		CaseType:            caseData.CaseType.Name,
		CaseID:              caseData.ID.String(),
		CreatedBy:           fmt.Sprintf("%s - %s", caseData.Creator.Name, caseData.Creator.Center.Name),
		CreatedDate:         caseData.CreatedAt.Format("2006-01-02 15:04:05"),
		VerifyStatus:        caseData.VerifyStatus,
		Channel:             caseData.Channel,
		Priority:            caseData.Priority,
		ReasonCode:          caseData.ReasonCode,
		DueDate:             caseData.DueDate.Format("2006-01-02"),
		AllocateToQueueTeam: utils.UUIDPtrToStringPtr(caseData.QueueID),
		CaseDescription:     caseData.Description,
		Status:              caseData.Status.Name,
		CurrentQueue:        caseData.Queue.Name,
	}

	return caseDetail, nil
}

func (uc *CaseUseCase) CreateCase(ctx *gin.Context, createdByID uuid.UUID, caseReq *model.CreateCaseRequest) (uuid.UUID, error) {
	statusMap, _ := uc.repo.LoadCaseStatus(ctx)
	caseTypeMap, _ := uc.repo.LoadCaseType(ctx)

	var caseTypeID uuid.UUID
	if caseReq.CaseTypeID != "" {
		id, err := utils.ParseUUID(caseReq.CaseTypeID)
		if err != nil {
			return uuid.Nil, err
		}
		caseTypeID = id
	}

	queueID, err := utils.ParseOptionalUUID(caseReq.AllocateToQueueTeam)
	if err != nil {
		return uuid.Nil, err
	}

	dispositionMainID, err := utils.ParseOptionalUUID(caseReq.DispositionMainID)
	if err != nil {
		return uuid.Nil, err
	}

	dispositionSubID, err := utils.ParseOptionalUUID(caseReq.DispositionSubID)
	if err != nil {
		return uuid.Nil, err
	}

	productID, err := utils.ParseOptionalUUID(caseReq.ProductID)
	if err != nil {
		return uuid.Nil, err
	}

	var priority string
	if caseTypeID == caseTypeMap["Inquiry and disposition"] {
		priority = "Normal"
	} else {
		priority = caseReq.Priority
	}

	var dueDate *time.Time
	if caseReq.DueDate != nil && *caseReq.DueDate != "" {
		parsedDueDate, err := time.Parse("2006-01-02", *caseReq.DueDate)
		if err != nil {
			return uuid.Nil, fmt.Errorf("invalid due date format: %w", err)
		}
		dueDate = &parsedDueDate
	}

	caseToSave := &model.Cases{
		CaseTypeID:        caseTypeID,
		CustomerName:      caseReq.CustomerName,
		CustomerID:        caseReq.CustomerID,
		DispositionMainID: dispositionMainID,
		DispositionSubID:  dispositionSubID,
		VerifyStatus:      caseReq.VerifyStatus,
		QueueID:           queueID,
		Channel:           caseReq.Channel,
		Description:       caseReq.CaseDescription,
		ReasonCode:        caseReq.ReasonCode,
		AssignedToUserID:  &createdByID,
		ProductID:         productID,
		Priority:          priority,
		StatusID:          statusMap["created"],
		StartDate:         time.Now(),
		EndDate:           time.Now().Add(72 * time.Hour),
		DueDate:           dueDate,
		CreatedBy:         createdByID,
		UpdatedBy:         createdByID,
	}

	caseId, err := uc.repo.CreateCaseInquiry(ctx, caseToSave)
	if err != nil {
		return uuid.Nil, err
	}

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
