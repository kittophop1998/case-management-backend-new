package database

import (
	"case-management/internal/domain/model"
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type CasePg struct {
	db *gorm.DB
}

func NewCasePg(db *gorm.DB) *CasePg {
	return &CasePg{db: db}
}

func (c *CasePg) GetAllCase(ctx context.Context, offset, limit int, filter model.CaseFilter, category string, currID uuid.UUID) ([]*model.Cases, int, error) {
	var cases []*model.Cases
	var total int64

	query := c.db.WithContext(ctx).Model(&model.Cases{}).
		Preload("Status").
		Preload("CaseType").
		Preload("Queue").
		Preload("Creator.Center").
		Preload("AssignedToUser.Center").
		Joins("LEFT JOIN cases_types as ct ON ct.id = cases.case_type_id").
		Joins("LEFT JOIN cases_status as cs ON cs.id = cases.status_id")

	// ---- category filter ----
	if category != "" {
		switch category {
		case "myCase":
			query = query.Where("assigned_to_user_id = ?", currID)
		case "availableCase":
			query = query.Where("assigned_to_user_id IS NULL")
		case "inquiryLog":
			query = query.Where("ct.name = ?", "Inquiry and disposition")
		case "caseHistory":
			query = query.Where("cs.name = ?", "closed")
		}
	}

	// ---- keyword ----
	if filter.Keyword != "" {
		query = query.Where("customer_name ILIKE ?", "%"+strings.TrimSpace(filter.Keyword)+"%")
	}

	// ---- priority ----
	if len(filter.Priority) > 0 {
		query = query.Where("priority IN ?", filter.Priority)
	}

	// ---- status ----
	if len(filter.StatusID) > 0 {
		query = query.Where("status_id IN ?", filter.StatusID)
	}

	// ---- queue ----
	if filter.QueueID != nil && *filter.QueueID != uuid.Nil {
		query = query.Where("queue_id = ?", *filter.QueueID)
	}

	// ---- นับ total ก่อน ----
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// ---- ดึง data ----
	if err := query.Offset(offset).Limit(limit).Order(filter.Sort).Find(&cases).Error; err != nil {
		return nil, 0, err
	}

	return cases, int(total), nil
}

func (c *CasePg) GetCaseByID(ctx context.Context, id uuid.UUID) (*model.Cases, error) {
	var cases model.Cases
	if err := c.db.WithContext(ctx).
		Preload("Status").
		Preload("ReasonCode").
		Preload("CaseType").
		Preload("Queue").
		Preload("AssignedToUser.Center").
		Preload("Creator.Center").
		First(&cases, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &cases, nil
}

func (c *CasePg) CreateCase(ctx context.Context, caseToSave *model.Cases) (uuid.UUID, error) {
	if err := c.db.WithContext(ctx).Create(caseToSave).Error; err != nil {
		return uuid.Nil, err
	}

	return caseToSave.ID, nil
}

func (c *CasePg) UpdateCaseDetail(ctx context.Context, caseToSave *model.Cases) error {
	if err := c.db.WithContext(ctx).Where("id = ?", caseToSave.ID).Updates(caseToSave).Error; err != nil {
		return err
	}
	return nil
}

func (c *CasePg) CreateCaseDispositionMains(ctx context.Context, input []*model.CaseDispositionMain) error {
	if err := c.db.WithContext(ctx).Create(input).Error; err != nil {
		return err
	}

	return nil
}

func (c *CasePg) CreateCaseDispositionSubs(ctx context.Context, input []*model.CaseDispositionSub) error {
	if err := c.db.WithContext(ctx).Create(input).Error; err != nil {
		return err
	}

	return nil
}

func (c *CasePg) GetCaseDispositionMains(ctx context.Context, caseID uuid.UUID) ([]*model.CaseDispositionMain, error) {
	var mains []*model.CaseDispositionMain
	if err := c.db.WithContext(ctx).
		Preload("Main").
		Where("case_id = ?", caseID).
		Find(&mains).Error; err != nil {
		return nil, err
	}

	return mains, nil
}

func (c *CasePg) GetCaseDispositionSubs(ctx context.Context, caseID uuid.UUID) ([]*model.CaseDispositionSub, error) {
	var subs []*model.CaseDispositionSub
	if err := c.db.WithContext(ctx).
		Preload("Sub").
		Where("case_id = ?", caseID).
		Find(&subs).Error; err != nil {
		return nil, err
	}

	return subs, nil
}

func (c *CasePg) CreateNoteType(ctx context.Context, note model.NoteTypes) (*model.NoteTypes, error) {
	if err := c.db.WithContext(ctx).Create(&note).Error; err != nil {
		return nil, err
	}
	return &note, nil
}

func (c *CasePg) GetNoteTypeByID(ctx context.Context, noteTypeID uuid.UUID) (*model.NoteTypes, error) {
	var noteType model.NoteTypes
	if err := c.db.WithContext(ctx).Where("id = ?", noteTypeID).First(&noteType).Error; err != nil {
		return nil, err
	}
	return &noteType, nil
}

func (r *CasePg) GetAllDisposition(ctx context.Context) ([]model.DispositionItem, error) {
	// ดึง DispositionMain ทั้งหมด พร้อม Subs
	var dispositionMains []model.DispositionMain
	if err := r.db.WithContext(ctx).
		Preload("Subs").
		Find(&dispositionMains).Error; err != nil {
		return nil, err
	}

	// Map ไปเป็น []DispositionItem (response)
	var dispositions []model.DispositionItem
	for _, main := range dispositionMains {
		item := model.DispositionItem{
			DispositionMain: model.DispositionMainRes{
				ID: main.ID.String(),
				TH: main.NameTH,
				EN: main.NameEN,
			},
		}

		for _, sub := range main.Subs {
			item.DispositionSubs = append(item.DispositionSubs, model.DispositionSubRes{
				ID: sub.ID.String(),
				TH: sub.NameTH,
				EN: sub.NameEN,
			})
		}

		dispositions = append(dispositions, item)
	}

	return dispositions, nil
}

func (r *CasePg) LoadCaseStatus(ctx context.Context) (map[string]uuid.UUID, error) {
	statusMap := make(map[string]uuid.UUID)
	var statuses []model.CaseStatus
	if err := r.db.WithContext(ctx).Find(&statuses).Error; err != nil {
		return nil, err
	}

	for _, status := range statuses {
		statusMap[status.Name] = status.ID
	}

	return statusMap, nil
}

func (r *CasePg) LoadCaseType(ctx context.Context) (map[string]uuid.UUID, error) {
	typeMap := make(map[string]uuid.UUID)
	var types []model.CaseTypes
	if err := r.db.WithContext(ctx).Find(&types).Error; err != nil {
		return nil, err
	}

	for _, t := range types {
		typeMap[t.Name] = t.ID
	}

	return typeMap, nil
}

func (c *CasePg) GetCaseNotes(ctx context.Context, caseID uuid.UUID) ([]*model.CaseNotes, error) {
	var notes []*model.CaseNotes
	if err := c.db.WithContext(ctx).
		Preload("Creator.Center").
		Where("case_id = ?", caseID).
		Order("created_at ASC").
		Find(&notes).Error; err != nil {
		return nil, err
	}

	return notes, nil
}

func (c *CasePg) AddCaseNote(ctx context.Context, note *model.CaseNotes) (uuid.UUID, error) {
	if err := c.db.WithContext(ctx).Create(note).Error; err != nil {
		return uuid.Nil, err
	}
	return note.ID, nil
}

func (r *CasePg) GenCaseCode(ctx context.Context) (string, error) {
	var lastCase model.Cases
	if err := r.db.WithContext(ctx).
		Order("created_at DESC").
		Limit(1).
		Debug().Find(&lastCase).Error; err != nil {
		return "", err
	}
	// ถ้าไม่มีเคสเลย ให้เริ่มที่ 1
	newNumber := 1

	if lastCase.Code != "" {
		parts := strings.Split(lastCase.Code, "_")
		if len(parts) == 2 {
			if n, err := strconv.Atoi(parts[1]); err == nil {
				newNumber = n + 1
			}
		}
	}

	// สร้างรหัสเคสใหม่ ในรูปแบบ CASE_001, CASE_002, ...
	newCode := fmt.Sprintf("CASE_%03d", newNumber)
	return newCode, nil
}

func (r *CasePg) GetCaseTypeByID(ctx context.Context, id uuid.UUID) (*model.CaseTypes, error) {
	var caseType model.CaseTypes
	if err := r.db.WithContext(ctx).First(&caseType, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &caseType, nil
}
