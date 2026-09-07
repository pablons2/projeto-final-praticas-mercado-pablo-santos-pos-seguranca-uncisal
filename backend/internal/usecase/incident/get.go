package incident

import (
	"context"

	"incidenttrack/internal/domain/entity"
)

// Execute detalha um incidente. Retorna errdomain.ErrNotFound se não
// existe e errdomain.ErrForbidden se pertence a outro usuário.
func (uc *GetIncident) Execute(ctx context.Context, id, ownerID string) (*entity.Incident, error) {
	return uc.repo.GetByIDForOwner(ctx, id, ownerID)
}
