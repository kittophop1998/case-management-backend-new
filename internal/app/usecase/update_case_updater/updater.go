package updatecaseupdater

import (
	"context"

	"github.com/google/uuid"
)

type CaseUpdater interface {
	Update(ctx context.Context, caseID uuid.UUID, data map[string]interface{}) error
}
