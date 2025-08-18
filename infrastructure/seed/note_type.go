package seed

import (
	"case-management/internal/domain/model"

	"gorm.io/gorm"
)

func SeedNoteTypes(db *gorm.DB) {
	noteTypes := []model.NoteTypes{
		{Name: "General", Description: "General notes"},
		{Name: "Follow-up", Description: "Follow-up notes"},
		{Name: "Urgent", Description: "Urgent notes"},
	}

	for _, noteType := range noteTypes {
		if err := db.FirstOrCreate(&noteType, model.NoteTypes{Name: noteType.Name}).Error; err != nil {
			panic("failed to seed note types: " + err.Error())
		}
	}
}
