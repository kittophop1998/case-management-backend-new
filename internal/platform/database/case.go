package database

import (
	"case-management/internal/domain/model"
	"encoding/json"
	"strings"

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

func (c *CasePg) CreateCaseInquiry(ctx *gin.Context, caseToSave *model.Cases) (uuid.UUID, error) {
	if err := c.db.FirstOrCreate(caseToSave).Error; err != nil {
		return uuid.Nil, err
	}

	return caseToSave.ID, nil
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

func (c *CasePg) GetAllCase(ctx *gin.Context, offset, limit int) ([]*model.Cases, int, error) {
	var cases []*model.Cases

	query := c.db.WithContext(ctx).Model(&model.Cases{}).
		Preload("Status")

	if err := query.Limit(limit).Offset(offset).Find(&cases).Error; err != nil {
		return nil, 0, err
	}

	return cases, 10, nil
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

func (r *CasePg) GetAllDisposition(ctx *gin.Context, limit, offset int) ([]model.DispositionMain, int, error) {
	query := r.db.WithContext(ctx).Model(&model.DispositionMain{})

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

func (r *CasePg) GetAllDispositionNew(ctx *gin.Context, filter model.DispositionFilter) ([]model.DispositionMain, error) {
	var mains []model.DispositionMain

	query := r.db.WithContext(ctx).
		Model(&model.DispositionMain{}).
		Preload("Subs")

	if filter.Keyword != "" {
		like := "%" + filter.Keyword + "%"

		query = query.Where(`
			disposition_mains.name ILIKE ? OR 
			disposition_mains.description ILIKE ? OR 
			id IN (
				SELECT main_id 
				FROM disposition_subs 
				WHERE name ILIKE ? OR description ILIKE ?
			)`,
			like, like, like, like)
	}

	if filter.Limit > 0 {
		query = query.Limit(filter.Limit)
	}

	if filter.Offset > 0 {
		query = query.Offset(filter.Offset)
	}

	if err := query.Find(&mains).Error; err != nil {
		return nil, err
	}

	return mains, nil
}

func (r *CasePg) LoadCaseStatus(ctx *gin.Context) (map[string]uuid.UUID, error) {
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
