package database

import (
	"case-management/internal/domain/model"
	"encoding/json"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type CasePg struct {
	db *gorm.DB
}

func NewCasePg(db *gorm.DB) *CasePg {
	return &CasePg{db: db}
}

func (c *CasePg) CreateCase(ctx *gin.Context, data *model.CreateCaseRequest) (uuid.UUID, error) {
	userIDStr := ctx.GetString("userId")
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return uuid.Nil, err
	}

	newToSave := &model.Cases{
		Model: model.Model{
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		CustomerId:        data.CustomerId,
		CaseTypeId:        data.CaseTypeId,
		DispositionMainId: data.DispositionMainId,
		Description:       data.CaseDescription,
		CreatedBy:         userID,
		UpdatedBy:         userID,
	}

	if err := c.db.Create(newToSave).Error; err != nil {
		return uuid.Nil, err
	}

	return newToSave.ID, nil
}

func (c *CasePg) CreateCaseDispositionMains(ctx *gin.Context, data datatypes.JSON) error {
	var dispositions []model.CaseDispositionMain
	if err := json.Unmarshal(data, &dispositions); err != nil {
		return err
	}

	for _, disposition := range dispositions {
		if err := c.db.WithContext(ctx).Create(&disposition).Error; err != nil {
			return err
		}
	}

	return nil
}

func (c *CasePg) CreateCaseDispositionSubs(ctx *gin.Context, data datatypes.JSON) error {
	var dispositions []model.CaseDispositionSub
	if err := json.Unmarshal(data, &dispositions); err != nil {
		return err
	}

	for _, disposition := range dispositions {
		if err := c.db.WithContext(ctx).Create(&disposition).Error; err != nil {
			return err
		}
	}

	return nil
}

func (c *CasePg) GetAllCase(ctx *gin.Context, limit, offset int, filter model.CaseFilter) ([]*model.Cases, error) {
	var cases []*model.Cases
	query := c.db.Model(&model.Cases{}).Offset(offset).Limit(limit)

	if filter.StatusID != nil {
		query = query.Where("status_id = ?", filter.StatusID)
	}

	if err := query.Find(&cases).Error; err != nil {
		return nil, err
	}

	return cases, nil
}

func (c *CasePg) CountWithFilter(ctx *gin.Context, filter model.CaseFilter) (int, error) {
	var count int64
	query := c.db.WithContext(ctx).Model(&model.Cases{})

	if filter.Keyword != "" {
		kw := "%" + strings.TrimSpace(filter.Keyword) + "%"
		query = query.Where(
			c.db.Where("title ILIKE ?", kw).
				Or("customer_id ILIKE ?", kw).
				Or("created_by ILIKE ?", kw).
				Or("CAST(sla_date AS TEXT) ILIKE ?", kw).
				Or("CAST(created_at AS TEXT) ILIKE ?", kw),
		)
	}

	if filter.StatusID != nil {
		query = query.Where("status_id = ?", *filter.StatusID)
	}

	if filter.PriorityID != nil {
		query = query.Where("priority_id = ?", *filter.PriorityID)
	}

	if filter.SLADateFrom != nil {
		query = query.Where("sla_date >= ?", *filter.SLADateFrom)
	}

	if filter.SLADateTo != nil {
		query = query.Where("sla_date <= ?", *filter.SLADateTo)
	}

	if err := query.Count(&count).Error; err != nil {
		return 0, err
	}

	return int(count), nil
}

func (c *CasePg) CreateNoteType(ctx *gin.Context, note model.NoteTypes) (*model.NoteTypes, error) {
	if err := c.db.WithContext(ctx).Create(&note).Error; err != nil {
		return nil, err
	}
	return &note, nil
}

func (c *CasePg) GetCaseByID(ctx *gin.Context, id uuid.UUID) (*model.Cases, error) {
	var cases model.Cases
	if err := c.db.WithContext(ctx).First(&cases, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &cases, nil
}

func (c *CasePg) AddInitialDescription(ctx *gin.Context, caseID uuid.UUID, newDescription string) error {
	var caseRecord struct {
		InitialDescriptions datatypes.JSON `gorm:"type:jsonb"`
	}

	err := c.db.WithContext(ctx).
		Model(&model.Cases{}).
		Select("initial_descriptions").
		Where("id = ?", caseID).
		Take(&caseRecord).Error
	if err != nil {
		return err
	}

	var descriptions []string
	if len(caseRecord.InitialDescriptions) > 0 {
		if err := json.Unmarshal(caseRecord.InitialDescriptions, &descriptions); err != nil {
			descriptions = []string{}
		}
	} else {
		descriptions = []string{}
	}

	descriptions = append(descriptions, newDescription)

	updatedJSON, err := json.Marshal(descriptions)
	if err != nil {
		return err
	}

	return c.db.WithContext(ctx).
		Model(&model.Cases{}).
		Where("id = ?", caseID).
		Update("initial_descriptions", datatypes.JSON(updatedJSON)).Error
}

func (c *CasePg) GetNoteTypeByID(ctx *gin.Context, noteTypeID uuid.UUID) (*model.NoteTypes, error) {
	var noteType model.NoteTypes
	if err := c.db.WithContext(ctx).Where("id = ?", noteTypeID).First(&noteType).Error; err != nil {
		return nil, err
	}
	return &noteType, nil
}

func (r *CasePg) GetAllDisposition(ctx *gin.Context, limit, offset int, keyword string) ([]model.DispositionMain, int, error) {
	query := r.db.WithContext(ctx).Model(&model.DispositionMain{})

	if keyword != "" {
		like := "%" + keyword + "%"
		query = query.Where("name ILIKE ? OR description ILIKE ?", like, like)
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	var mains []model.DispositionMain
	if err := query.
		Limit(limit).
		Offset(offset).
		Order("name ASC").
		Find(&mains).Error; err != nil {
		return nil, 0, err
	}

	return mains, int(total), nil
}
