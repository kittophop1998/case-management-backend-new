package seed

import (
	"case-management/internal/domain/model"

	"gorm.io/gorm"
)

func SeedNoteTypes(db *gorm.DB) {
	noteTypes := []model.NoteTypes{
		{Name: "การติดตาม(Follow-Up)", Description: "การติดตาม(Follow-Up)"},
		{Name: "หมายเหตุพิเศษ", Description: "หมายเหตุพิเศษ"},
		{Name: "การร้องเรียน / ข้อเสนอแนะ", Description: "การร้องเรียน / ข้อเสนอแนะ"},
	}

	for _, noteType := range noteTypes {
		if err := db.FirstOrCreate(&noteType, model.NoteTypes{Name: noteType.Name}).Error; err != nil {
			panic("failed to seed note types: " + err.Error())
		}
	}
}
