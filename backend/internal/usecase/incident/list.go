package incident

import (
	"context"

	"incidenttrack/internal/domain/entity"
	"incidenttrack/internal/domain/repository"
)

// Execute devolve os incidentes de ownerID aplicando os filtros.
func (uc *ListIncidents) Execute(ctx context.Context, ownerID string, f repository.IncidentFilter) ([]entity.Incident, error) {
	if f.Limit <= 0 || f.Limit > 500 {
		f.Limit = 500
	}
	return uc.repo.ListByOwner(ctx, ownerID, f)
}
