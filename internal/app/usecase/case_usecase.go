package usecase

import (
	"case-management/internal/domain/model"
	"case-management/internal/domain/repository"
	"case-management/utils"
	"context"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
)

type CaseUseCase struct {
	repo repository.CaseRepository
}

func NewCaseUseCase(repo repository.CaseRepository) *CaseUseCase {
	return &CaseUseCase{repo: repo}
}

func (uc *CaseUseCase) GetAllCases(ctx context.Context, page, limit int, filter model.CaseFilter, category string, currID uuid.UUID) ([]*model.CaseResponse, int, error) {
	offset := (page - 1) * limit
	caseRepo, total, err := uc.repo.GetAllCase(ctx, offset, limit, filter, category, currID)
	if err != nil {
		return nil, 0, err
	}

	caseResponses := make([]*model.CaseResponse, 0, len(caseRepo))
	for _, c := range caseRepo {
		caseResponses = append(caseResponses, &model.CaseResponse{
			Code:         c.Code,
			CustomerID:   c.CustomerID,
			CustomerName: c.CustomerName,
			Status:       c.Status.Name,
			CaseType:     c.CaseType.Name,
			CurrentQueue: c.Queue.Name,
			CurrentUser:  utils.UserNameCenter(*c.AssignedToUser),
			SLADate:      utils.FormatDate(&c.EndDate, "2006-01-02 15:04"),
			CaseID:       c.ID.String(),
			AeonID:       c.AeonID,
			CaseGroup:    c.CaseType.Group,
			CreatedBy:    utils.UserNameCenter(c.Creator),
			CreatedDate:  utils.FormatDate(&c.CreatedAt, "2006-01-02 15:04"),
			CasePriority: c.Priority,
			ClosedDate:   utils.FormatDate(&c.ClosedDate, "2006-01-02 15:04"),
			ReceivedFrom: "Fraud",
		})
	}
	return caseResponses, total, nil

}

func (uc *CaseUseCase) GetCaseByID(ctx context.Context, id uuid.UUID) (*model.CaseDetailResponse, error) {
	caseData, err := uc.repo.GetCaseByID(ctx, id)
	if err != nil {
		return nil, err
	}

	caseDetail := &model.CaseDetailResponse{
		Code:                caseData.Code,
		CaseType:            caseData.CaseType.Name,
		CaseTypeID:          caseData.CaseType.ID.String(),
		CaseGroup:           caseData.CaseType.Group,
		CaseID:              caseData.ID.String(),
		CreatedBy:           utils.UserNameCenter(caseData.Creator),
		VerifyStatus:        caseData.VerifyStatus,
		Channel:             caseData.Channel,
		Priority:            caseData.Priority,
		ReasonCode:          utils.UUIDPtrToStringPtr(caseData.ReasonCodeID),
		AllocateToQueueTeam: utils.UUIDPtrToStringPtr(caseData.QueueID),
		CaseDescription:     caseData.Description,
		Status:              caseData.Status.Name,
		CurrentQueue:        caseData.Queue.Name,
		CreatedDate:         utils.FormatDate(&caseData.CreatedAt, "2006-01-02 15:04"),
		DueDate:             utils.FormatDate(caseData.DueDate, "2006-01-02"),
	}

	return caseDetail, nil
}

func (uc *CaseUseCase) CreateCase(ctx context.Context, createdByID uuid.UUID, caseReq *model.CreateCaseRequest) (uuid.UUID, error) {
	statusMap, _ := uc.repo.LoadCaseStatus(ctx)

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

	reasonCodeID, err := utils.ParseOptionalUUID(caseReq.ReasonCode)
	if err != nil {
		return uuid.Nil, err
	}

	dueDate, err := utils.ParseOptionalDate(caseReq.DueDate, "2006-01-02")
	if err != nil {
		return uuid.Nil, fmt.Errorf("invalid dueDate format: %v", err)
	}

	code, err := uc.repo.GenCaseCode(ctx)
	if err != nil {
		log.Fatal(err)
	}

	caseToSave := &model.Cases{
		Code:              code,
		CaseTypeID:        caseTypeID,
		CustomerName:      caseReq.CustomerName,
		CustomerID:        caseReq.CustomerID,
		DispositionMainID: dispositionMainID,
		DispositionSubID:  dispositionSubID,
		VerifyStatus:      caseReq.VerifyStatus,
		QueueID:           queueID,
		Channel:           caseReq.Channel,
		Description:       caseReq.CaseDescription,
		ReasonCodeID:      reasonCodeID,
		AssignedToUserID:  &createdByID,
		ProductID:         productID,
		Priority:          caseReq.Priority,
		StatusID:          statusMap["created"],
		StartDate:         time.Now(),
		EndDate:           time.Now().Add(72 * time.Hour),
		DueDate:           dueDate,
		CreatedBy:         createdByID,
		UpdatedBy:         createdByID,
	}

	caseId, err := uc.repo.CreateCase(ctx, caseToSave)
	if err != nil {
		return uuid.Nil, err
	}

	return caseId, nil
}

func (uc *CaseUseCase) UpdateCaseDetail(ctx context.Context, caseID uuid.UUID, input *model.UpdateCaseRequest) error {
	ReasonCodeID, err := utils.ParseOptionalUUID(&input.ReasonCodeID)
	if err != nil {
		return err
	}

	dueDate, err := utils.ParseOptionalDate(&input.DueDate, "2006-01-02")
	if err != nil {
		return err
	}

	queueID, err := utils.ParseOptionalUUID(&input.ReallocateToQueueTeam)
	if err != nil {
		return err
	}

	caseToSave := &model.Cases{
		ID:           caseID,
		Priority:     input.Priority,
		ReasonCodeID: ReasonCodeID,
		DueDate:      dueDate,
		QueueID:      queueID,
	}
	return uc.repo.UpdateCaseDetail(ctx, caseToSave)
}

func (uc *CaseUseCase) GetCaseNotes(ctx context.Context, caseID uuid.UUID) ([]*model.CaseNotesResponse, error) {
	notes, err := uc.repo.GetCaseNotes(ctx, caseID)
	if err != nil {
		return nil, err
	}

	var responses []*model.CaseNotesResponse
	for _, note := range notes {
		responses = append(responses, &model.CaseNotesResponse{
			ID:        note.ID,
			Content:   note.Content,
			CreatedBy: utils.UserNameCenter(note.Creator),
			CreatedAt: note.CreatedAt,
		})
	}

	return responses, nil
}

func (uc *CaseUseCase) AddCaseNote(ctx context.Context, createdByID uuid.UUID, caseID uuid.UUID, input *model.CaseNoteRequest) (uuid.UUID, error) {
	note := &model.CaseNotes{
		CaseId:    caseID,
		CreatedBy: createdByID,
		Content:   input.Content,
	}

	return uc.repo.AddCaseNote(ctx, note)
}

func (uc *CaseUseCase) AddInitialDescription(c context.Context, caseID string, newDescription string) error {
	caseUUID, err := uuid.Parse(caseID)
	if err != nil {
		return err
	}
	return uc.repo.AddInitialDescription(c, caseUUID, newDescription)
}

func (uc *CaseUseCase) GetAllDisposition(ctx context.Context) ([]model.DispositionItem, error) {
	return uc.repo.GetAllDisposition(ctx)
}

func (uc *CaseUseCase) GetCaseTypeByID(ctx context.Context, caseTypeID uuid.UUID) (model.CaseTypes, error) {
	caseType, err := uc.repo.GetCaseTypeByID(ctx, caseTypeID)
	if err != nil {
		return model.CaseTypes{}, err
	}

	return *caseType, nil
}
